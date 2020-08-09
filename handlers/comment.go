package handlers

import (
	"net/http"
	"encoding/json"
	"exp/models"
	"github.com/gorilla/mux"
	"strconv"
)

// GetComment returns a post from the database
func GetComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID, id := params["postID"], params["id"]
	var comment models.Comment
	db.First(&comment, "PostID = ? AND id = ?", postID, id)
	json.NewEncoder(w).Encode(&comment)
}

// CreateComment adds new comment 
func CreateComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["postID"]
	var (
	comment models.Comment
	)
	json.NewDecoder(r.Body).Decode(&comment)
	temp, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		json.NewEncoder(w).Encode("Invalid Post ID!")
	}
	if temp == uint64(comment.PostID) {
		db.Create(&comment)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&comment)
	} else {
		json.NewEncoder(w).Encode(models.ErrorResponse{Err : "Mismatch Post ID!"})
	}
}

// FetchComments returns the comments of a post
func FetchComments(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["postID"]
	var comments []models.Comment
	db.Find(&comments, "PostID = ?", postID)
	json.NewEncoder(w).Encode(&comments)
}

// UpdateComment modifies a comment
func UpdateComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID, id := params["postID"], params["id"]
	var comment models.Comment
	db.First(&comment, "PostID = ? AND id = ?", postID, id)
	json.NewDecoder(r.Body).Decode(&comment)
	db.Save(&comment)
	json.NewEncoder(w).Encode(&comment)
}

// DeleteComment removes a comment
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID, id := params["postID"], params["id"]
	var comment models.Comment
	db.First(&comment, "PostID = ? AND id = ?", postID, id)
	db.Delete(&comment)
	json.NewEncoder(w).Encode("Comment deleted")
}