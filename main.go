package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

//Date Struct
type Date struct {
	Day   int
	Month int
	Year  int
}

//Article Struct
type Article struct {
	Title      string
	Body       string
	UploadDate Date
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/posts", PostsHandler)

	http.ListenAndServe(":5000", r)
}

//PostsHandler return all posts
func PostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode({"How To Sing", "Be Patient and Keep Calm, Don't Forget to Brush your Teeth.", Date{28, 8, 2002}})
}


