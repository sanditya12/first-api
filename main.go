package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	. "github.com/sanditya12/rest-api/constants"
)

//Date Struct
type Date struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

//Post Struct
type Post struct {
	Title      string `json:"title"`
	Body       string `json:"body"`
	UploadDate Date   `json:"uploadDate"`
}

var posts []Post

var db *sql.DB

func main() {
	pgURL, err := pq.ParseURL(EsqlURL)
	checkErr(err)

	db, err = sql.Open("postgres", pgURL)
	checkErr(err)

	err = db.Ping()
	checkErr(err)

	r := mux.NewRouter()
	r.HandleFunc("/api/posts", PostsHandler)
	r.HandleFunc("/api/posts/{id}", postHandler)
	r.HandleFunc("/api/add", addPost).Methods("POST")

	log.Fatal(http.ListenAndServe(":5000", r))
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

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
