package handlers

import (
	"exp/models"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
)

// FetchPosts returns a slice of posts in the database
func FetchPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	db.Find(&posts)
	json.NewEncoder(w).Encode(&posts)
}

// UpdatePost modifies a post
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	params := mux.Vars(r)
	id := params["id"]
	db.First(&post, id)
	json.NewDecoder(r.Body).Decode(&post)
	db.Save(&post)
	json.NewEncoder(w).Encode(&post)
}

// DeletePost removes a post in the database
func DeletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var post models.Post
	db.First(&post, id)
	db.Delete(&post)
	json.NewEncoder(w).Encode("Post deleted!")
}

// CreatePost add a new post to the database
func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	json.NewDecoder(r.Body).Decode(&post)
	db.Create(&post)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&post)
}

// GetPost returns a post from the database
func GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var post models.Post
	db.First(&post, id)
	json.NewEncoder(w).Encode(&post)
}