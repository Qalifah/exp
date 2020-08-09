package handlers

import (
	"exp/models"
	"exp/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"os"
	"log"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"github.com/joho/godotenv"
	"github.com/twinj/uuid"
	gocache "github.com/patrickmn/go-cache"
)

type error interface {
	Error() string
}

var (
	db = utils.ConnectDB()
	cache = gocache.New(time.Minute * 5, time.Second * 30)
)

func init() {
	var e = godotenv.Load()
	if e != nil {
		log.Fatal("Error loading .env file")
	}
}

// HelloHandler is a test handler
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Wagwan Gee!!")
}

// TestHandler is a test handler
func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("let's get started!"))
}

//Login handles identifying the user accessing the site
func Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		json.NewEncoder(w).Encode("Invalid request")
		return
	}
	resp := Authenticate(user.Email, user.Password)
	json.NewEncoder(w).Encode(resp)
}

//Authenticate verifies the user
func Authenticate(email, password string) map[string]interface{} {
	var user models.User

	if err := db.Where("Email = ?", email).First(&user).Error; err != nil {
		var resp = map[string]interface{}{"status":false, "message": "Email address not found!"}
		return resp
	}
	
	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword {
		resp := map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again!"}
		return resp
	}
	expiresAt := time.Now().Add(time.Minute * 15).Unix()
	cc := &models.CustomClaim {
		UserID : user.ID,
		Name : user.Username,
		Email : user.Email,
		StandardClaims : &jwt.StandardClaims {
			ExpiresAt : expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cc)
	tokenString, err := token.SignedString([]byte(os.Getenv("SIGNING_KEY")))
	if err != nil {
		fmt.Println(err)
	}

	refExpiresAt := time.Now().Add(time.Hour * 24 * 3).Unix()
	refCC := &models.CustomClaim {
		UserID : user.ID,
		Name : user.Username,
		Email : user.Email,
		StandardClaims : &jwt.StandardClaims {
			ExpiresAt : refExpiresAt,
		},
	}
	refToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refCC)
	refTokenString, err := refToken.SignedString([]byte(os.Getenv("REFRESH_TOKEN_KEY")))
	if err != nil {
		fmt.Println(err)
	}
	
	resp := map[string]interface{}{"message": "logged in!"}
	resp["accessToken"] = tokenString
	resp["refreshToken"] = refTokenString
	resp["user"] = user
	return resp
}

//CreateToken generates new token for authentication/authorization
func CreateToken(userID uint) (map[string]string, error) {
	expiresAt := time.Now().Add(time.Minute * 15).Unix()
	accessUUID := uuid.NewV4().String()
	cc := jwt.MapClaims{}
	cc["UserID"] = userID
	cc["access_uuid"] = accessUUID
	cc["ExpiresAt"] = expiresAt

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cc)
	tokenString, err := token.SignedString([]byte(os.Getenv("SIGNING_KEY")))
	if err != nil {
		return nil, err
	}

	refExpiresAt := time.Now().Add(time.Hour * 24 * 3).Unix()
	tokenUUID := uuid.NewV4().String()
	refCC := jwt.MapClaims{}
	refCC["UserID"] = userID
	refCC["token_uuid"] = tokenUUID
	refCC["ExpiresAt"] = refExpiresAt
	refToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refCC)
	refTokenString, err := refToken.SignedString([]byte(os.Getenv("REFRESH_TOKEN_KEY")))
	if err != nil {
		return nil, err
	}
	
	resp := map[string]string{}
	resp["accessToken"] = tokenString
	resp["refreshToken"] = refTokenString
	return resp, nil
}

// Register initiate the user
func Register(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		err := models.ErrorResponse{
			Err : "Password encryption failed",
		}
		json.NewEncoder(w).Encode(&err)
	}
	user.Password = string(pass)

	createdUser := db.Create(user)
	var errMessage = createdUser.Error

	if createdUser.Error != nil {
		fmt.Println(errMessage)
	}
	json.NewEncoder(w).Encode(createdUser)
}

// Logout does what it always do
func Logout(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Authorization")
	json.NewEncoder(w).Encode(map[string]interface{}{"message" : "You're logged out!"})
}

// FetchUsers returns users in the database
func FetchUsers(w http.ResponseWriter, r *http.Request) {
	users := &[]models.User{}
	if val, found := cache.Get("allUsers"); found {
		users = val.(*[]models.User)
	} else {
		db.Find(users)
		cache.Set("allUsers", users, gocache.DefaultExpiration)
	}
	json.NewEncoder(w).Encode(users)
}

//UpdateUser modifies a user 
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	params := mux.Vars(r)
	id := params["id"]
	db.First(user, id)
	json.NewDecoder(r.Body).Decode(user)
	db.Save(user)
	json.NewEncoder(w).Encode(user)
}

// DeleteUser removes a user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var user models.User
	db.First(&user, id)
	db.Delete(&user)
	json.NewEncoder(w).Encode("User deleted")
}

//GetUser returns a particular user
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var user models.User
	db.First(&user, id)
	json.NewEncoder(w).Encode(&user)
}
