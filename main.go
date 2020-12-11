package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//Date Struct
type Date struct {
	Day   int
	Month int
	Year  int
}

//Post Struct
type Post struct {
	Title      string
	Body       string
	UploadDate Date
}

var posts []Post

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/posts", PostsHandler)
	r.HandleFunc("/api/posts/{id}", postHandler)
	r.HandleFunc("/api/add", addPost).Methods("POST")

	http.ListenAndServe(":5000", r)
}

//PostsHandler return all posts
func PostsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(posts)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	// fmt.Printf("%T",id)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(posts[id])
}

func addPost(w http.ResponseWriter, r *http.Request) {
	var newPost Post
	json.NewDecoder(r.Body).Decode(&newPost)

	posts = append(posts, newPost)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// {"How To Sing", "Be Patient and Keep Calm, Don't Forget to Brush your Teeth.", Date{28, 8, 2002}}
