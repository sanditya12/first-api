package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	. "github.com/sanditya12/rest-api/constants"
	"golang.org/x/crypto/bcrypt"
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
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//JWT struct
type JWT struct {
	Token string `json:"token"`
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
	r.HandleFunc("/api/login", login).Methods("POST")

	log.Fatal(http.ListenAndServe(":5000", r))
}

func signup(w http.ResponseWriter, r *http.Request) {
	var user User
	var error Error
	json.NewDecoder(r.Body).Decode(&user)
	if user.Email == "" {
		error.Message = "Email is Missing"
		resError(w, http.StatusBadRequest, error)
		return
	}
	if user.Password == "" {
		error.Message = "Password is Missing"
		resError(w, http.StatusBadRequest, error)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	checkErr(err)

	user.Password = string(hashed)

	queryBody := "insert into users (email, password) values($1,$2) RETURNING id;"

	err = db.QueryRow(queryBody, user.Email, user.Password).Scan(&user.ID)
	checkErr(err)
	user.Password = ""

	resJSON(w, user)
}

func login(w http.ResponseWriter, r *http.Request) {
	var user User
	var error Error
	var jwt JWT

	json.NewDecoder(r.Body).Decode(&user)

	if user.Email == "" {
		error.Message = "Email is Missing"
		resError(w, http.StatusBadRequest, error)
		return
	}
	if user.Password == "" {
		error.Message = "Password is Missing"
		resError(w, http.StatusBadRequest, error)
		return
	}
	password := user.Password
	queryBody := "select * from users where email=$1"
	err := db.QueryRow(queryBody, user.Email).Scan(&user.ID, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		error.Message = "This Email is Not Registered Yet"
		resError(w, http.StatusBadRequest, error)
		return
	}
	checkErr(err)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		error.Message = "Wrong Password"
		resError(w, http.StatusBadRequest, error)
		return
	}
	token, err := createToken(user)
	checkErr(err)
	jwt.Token = token

	resJSON(w, jwt)
}

func createToken(user User) (string, error) {
	var err error
	secret := SECRET

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "course",
	})

	tokenString, err := token.SignedString([]byte(secret))
	checkErr(err)

	return tokenString, nil
}

//PostsHandler return all posts
func PostsHandler(w http.ResponseWriter, r *http.Request) {

	resJSON(w, posts)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	checkErr(err)
	w.Header().Set("Content-Type", "application/json")

	resJSON(w, posts[id])
}

func addPost(w http.ResponseWriter, r *http.Request) {
	var newPost Post
	json.NewDecoder(r.Body).Decode(&newPost)

	posts = append(posts, newPost)

	resJSON(w, posts)
}

func resJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func resError(w http.ResponseWriter, status int, err Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
