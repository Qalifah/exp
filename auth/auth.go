package auth

import (
	"exp/models"
	"exp/handlers"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"log"
	"fmt"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)
//Exception formats error
type Exception models.Exception

//KeyCtx is used to define a context key
type KeyCtx string

func init() {
	var err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	/*var  client *redis.Client
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
	   Addr: dsn, 
	})
	fmt.Println("Redis : ", client)
*/
}

// JWTVerify is a middleware that check the authenticity of a request
func JWTVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenstr := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "",-1)
		if tokenstr == ""  {
			ResponseAuthError(w, "Invalid token")
			return
		}
		
		var cc = &models.CustomClaim{}
		token, err := jwt.ParseWithClaims(tokenstr, cc, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SIGNING_KEY")), nil
		})

		if err != nil {
			ResponseAuthError(w, err.Error())
			return
		}
		if _, ok := token.Claims.(*models.CustomClaim); ok && token.Valid {
			ctx := context.WithValue(r.Context(), KeyCtx("user"), cc)
			next.ServeHTTP(w, r.WithContext(ctx))  
		} else {
			ResponseAuthError(w, "Couldn't handle this token!")
			return
		}
	})
}

//ResponseAuthError is a function that return authentication error as json
func ResponseAuthError(w http.ResponseWriter, str string) {
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(Exception{Message: str})
}

//Refresh function provides new access-token and refresh-token
func Refresh(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(KeyCtx("user")).(*models.CustomClaim).UserID
	tokenstr := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "",-1)
		if tokenstr == "" {
			ResponseAuthError(w, "Invalid token")
			return
		}
	var refCC = &models.CustomClaim{}
	token, err := jwt.ParseWithClaims(tokenstr, refCC, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_TOKEN_KEY")), nil
	})
	if err != nil {
		ResponseAuthError(w, err.Error())
		return
	}
	if _, ok := token.Claims.(*models.CustomClaim); !ok && !token.Valid {
		ResponseAuthError(w, "Invalid token!")
		return
	}
	result, err := handlers.CreateToken(userID)
	if err != nil {
		ResponseAuthError(w, err.Error())
		return
	}
	json.NewEncoder(w).Encode(result)
}