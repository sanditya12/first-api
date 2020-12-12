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
// type Date struct {
// 	Day   int `json:"day"`
// 	Month int `json:"month"`
// 	Year  int `json:"year"`
// }

//Post Struct
type Post struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	// UploadDate Date   `json:"uploadDate"`
}

//User Struct
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Error Struct
type Error struct {
	Message string `json:"message"`
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
	r.HandleFunc("/api/signup", signup).Methods("POST")

	log.Fatal(http.ListenAndServe(":5000", r))
}

func signup(w http.ResponseWriter, r *http.Request) {
	var user User
	var error Error
	json.NewDecoder(r.Body).Decode(&user)
	if user.Email == "" {
		error.Message = "Email is Missing"
		respondWithError(w, error)
	}
	if user.Password == "" {
		error.Message = "Password is Missing"
		respondWithError(w, error)
	}
}

//PostsHandler return all posts
func PostsHandler(w http.ResponseWriter, r *http.Request) {

	// w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(posts)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	checkErr(err)
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

func respondWithError(w http.ResponseWriter, error Error) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(error)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
