package routes

import (
	"exp/handlers"
	"exp/auth"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var cor = cors.New(cors.Options{
	AllowedMethods : []string{http.MethodGet, http.MethodDelete, http.MethodHead, http.MethodPatch, http.MethodOptions, http.MethodPost, http.MethodPut},
	AllowCredentials : true,
})

// Routers contain the routing system
func Routers() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Use(CommonMiddleware)

	api := router.PathPrefix("/v1").Subrouter()
	api.Use(auth.JWTVerify)

	//r.HandleFunc("/", handlers.HelloHandler).Methods("GET")
	//r.HandleFunc("/api", handlers.TestHandler).Methods("GET")
	aut := router.PathPrefix("/auth").Subrouter()

	// swagger:route POST /register users register
	// ---
	// Signing up as a user.
	// responses:
	//   200:
	//   	"$ref": "#/login"
	aut.HandleFunc("/register", handlers.Register).Methods("POST")

	// swagger:route POST /login user login
	// ---
	// Verify if the user is known or not.
	// responses:
	//   200:
	//   	"$ref": "#/"
	aut.HandleFunc("/login", handlers.Login).Methods("POST")

	// swagger:route POST /logout user logout
	// ---
	// Let's log the user out!.
	// responses:
	//   200:
	//   	"$ref": "#/login"
	aut.HandleFunc("/logout", handlers.Logout).Methods("POST")
	
	api.HandleFunc("/token/refresh", auth.Refresh).Methods("GET")

	// swagger:operation GET user getUser
	// ---
	// returns a particular user whose id matches the one of passed in the request.
	// parameters:
	// - name : id
	// decsription : id of a user
	// type : int
	// required : true
	// responses:
	// 		"200" : user
	api.HandleFunc("/user/{id}", handlers.GetUser).Methods("GET")

	// swagger:operation PUT user updateUser
	// ---
	// modify a particular user's attribute whose id matches the one of passed in the request.
	// parameters:
	// - name : id
	// decsription : id of a user
	// type : int
	// required : true
	// responses:
	// 		"200" : user
	api.HandleFunc("/user/{id}", handlers.UpdateUser).Methods("PUT")

	// swagger:operation DELETE user deleteUser
	// ---
	// deletes a particular user whose id matches the one of passed in the request.
	// parameters:
	// - name : id
	// decsription : id of a user
	// type : int
	// required : true
	// responses:
	// 		"200" : string
	api.HandleFunc("/user/{id}", handlers.DeleteUser).Methods("DELETE")

	// swagger:route GET user fetchUsers
	// ---
	// returns all users available.
	// responses:
	// 		"200" : []user
	api.HandleFunc("/users", handlers.FetchUsers).Methods("GET")

	// swagger:operation GET post user getPost
	// ---
	// returns a particular post whose id matches the one of passed in the request.
	// parameters:
	// - name : id
	// decsription : id of a post
	// type : int
	// required : true
	// responses:
	// 		"200" : post
	api.HandleFunc("/post/{id}", handlers.GetPost).Methods("GET")

	// swagger:operation UPDATE post user updatePost
	// ---
	// modify a particular post whose id matches the one of passed in the request.
	// parameters:
	// - name : id
	// decsription : id of a post
	// type : int
	// required : true
	// responses:
	// 		"200" : user
	api.HandleFunc("/post/{id}", handlers.UpdatePost).Methods("PUT")

	// swagger:operation DELETE post user deletePost
	// ---
	// deletes a particular post whose id matches the one of passed in the request.
	// parameters:
	// - name : id
	// decsription : id of a post
	// type : int
	// required : true
	// responses:
	// 		"200" : string
	api.HandleFunc("/post/{id}", handlers.DeletePost).Methods("DELETE")

	// swagger:route POST post user createPost
	// ---
	// create a new post.
	// responses:
	// 		"200" : post
	api.HandleFunc("/post", handlers.CreatePost).Methods("POST")

	// swagger:route GET post user fetchPosts
	// ---
	// returns all posts available.
	// responses:
	// 		"200" : []post
	api.HandleFunc("/posts", handlers.FetchPosts).Methods("GET")

	// swagger:operation GET post user comment getComment
	// ---
	// returns a particular comment whose id and postID matches the one of passed in the request.
	// parameters:
	// - name : postID
	// description : id of the post in which the comment was made
	// type : int
	// required : true
	// - name : id
	// description : id of a comment
	// type : int
	// required : true
	// responses:
	// 		"200" : comment
	api.HandleFunc("/post/{postID}/comment/{id}", handlers.GetComment).Methods("GET")

	// swagger:operation PUT post user comment updateComment
	// ---
	// modifies a particular comment whose id and postID matches the one of passed in the request.
	// parameters:
	// - name : postID
	// description : id of the post in which the comment was made
	// type : int
	// required : true
	// - name : id
	// description : id of a comment
	// type : int
	// required : true
	// responses:
	// 		"200" : comment
	api.HandleFunc("/post/{postID}/comment/{id}", handlers.UpdateComment).Methods("PUT")

	// swagger:operation DELETE post user comment deleteComment
	// ---
	// deletes a particular comment whose id and postID matches the one of passed in the request.
	// parameters:
	// - name : postID
	// description : id of the post in which the comment was made
	// type : int
	// required : true
	// - name : id
	// description : id of a comment
	// type : int
	// required : true
	// responses:
	// 		"200" : comment
	api.HandleFunc("/post/{postID}/comment/{id}", handlers.DeleteComment).Methods("DELETE")

	// swagger:operation POST post user comment createComment
	// ---
	// create new comment.
	// parameters:
	// - name : postID
	// description : id of the post in which the comment was made
	// type : int
	// required : true
	// responses:
	// 		"200" : post
	api.HandleFunc("post/{postID}/comment", handlers.CreateComment).Methods("POST")

	// swagger:operation GET post user comment fetchComment
	// ---
	// returns all comments of a post.
	// parameters:
	// - name : postID
	// description : id of the post in which the comment was made
	// type : int
	// required : true
	// responses:
	// 		"200" : post
	api.HandleFunc("post/{postID}/comments", handlers.FetchComments).Methods("GET")

	api.HandleFunc("/", handlers.TestHandler).Methods("GET")
	return router
}

// CommonMiddleware : name should be self explanatory, right?
func CommonMiddleware(next http.Handler) http.Handler {
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next = cor.Handler(next)
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}