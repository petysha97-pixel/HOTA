package main

import (
	"HOTA/internal/handlers"
	"HOTA/internal/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func main() {

	db, err := sql.Open("sqlite", "../db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	models.UserDB = db
	fmt.Println("Connected to SQLite")

	http.HandleFunc("/", handlers.Resistr)
	http.HandleFunc("/user", handlers.NewUser)
	http.HandleFunc("/profile/{id}", handlers.HOME)
	http.HandleFunc("/search", handlers.SearchUsers)

	// http.HandleFunc("/profile", profile.Profile)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

// sql.Open - db, err := sql.Open("sqlite", "app.db") Открывает БД.

// db.Exec Выполняет запрос без результата.
// _, err := db.Exec(`
// CREATE TABLE users(
//     id INTEGER PRIMARY KEY,
//     name TEXT
// )
// `)

// db.Query  Возвращает много строк.
// rows, err := db.Query(
//     "SELECT id, name FROM users",
// )

// db.QueryRow Возвращает одну строку.
// var name string

// err := db.QueryRow(
//     "SELECT name FROM users WHERE id = ?",
//     1,
// ).Scan(&name)

// rows.Next - Переход к следующей записи.
// for rows.Next() {
// }

// rows.Scan - Чтение данных.
// var id int
// var name string

// rows.Scan(&id, &name)

// db.Prepare - Подготовленный запрос.
// stmt, err := db.Prepare(
//     "INSERT INTO users(name) VALUES(?)",
// )

// stmt.Exec("John")
// stmt.Exec("Kate")
// stmt.Exec("Alex")
