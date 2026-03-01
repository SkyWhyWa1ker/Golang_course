package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Movie struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

var db *sql.DB

func connectDB() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	fmt.Println("Connected to database")
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title FROM movies")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var movies []Movie

	for rows.Next() {
		var m Movie
		rows.Scan(&m.ID, &m.Title)
		movies = append(movies, m)
	}

	json.NewEncoder(w).Encode(movies)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	var m Movie
	json.NewDecoder(r.Body).Decode(&m)

	_, err := db.Exec("INSERT INTO movies(title) VALUES($1)", m.Title)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func main() {

	connectDB()

	http.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			getMovies(w, r)
		}

		if r.Method == http.MethodPost {
			createMovie(w, r)
		}

	})

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
