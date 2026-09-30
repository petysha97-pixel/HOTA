package main

import (
	"HOTA/internal/handlers"
	projectHAND "HOTA/internal/handlers/project_handlers"
	slotHAND "HOTA/internal/handlers/project_handlers/slot_handlers"
	"HOTA/internal/models"
	"HOTA/internal/service"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Ошибка в загрузке файла .env")
	}
	fmt.Println("Подключен файл .env")

	db, err := sql.Open("sqlite", "../db.db?_pragma=busy_timeout(10000)&_pragma=foreign_keys(ON)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=temp_store(MEMORY)")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	models.UserDB = db
	fmt.Println("SQLite подключен")

	mux := http.NewServeMux()

	// Короткая обёртка над JWTMiddleware, чтобы не писать service.JWTMiddleware(http.HandlerFunc(...))
	auth := func(h http.HandlerFunc) http.Handler { return service.JWTMiddleware(h) }

	// ===== Регистрация / авторизация =====
	mux.HandleFunc("POST /user", handlers.NewUser)
	mux.HandleFunc("POST /user/auth", handlers.Auth)
	mux.HandleFunc("GET /stack", handlers.GetStacks)

	// ===== Профиль =====
	mux.Handle("GET /profile", auth(handlers.GetUser))
	mux.HandleFunc("GET /profile/other/{id}", handlers.GetUserOtherProfile)
	mux.HandleFunc("GET /users", handlers.GetAllUser)
	mux.HandleFunc("GET /users/searche", handlers.SearcheUsers)

	// ===== CRUD пользователя =====
	mux.Handle("PUT /user/update", auth(handlers.UpdateUser))
	mux.Handle("PATCH /user/update", auth(handlers.UpdateUser))
	mux.Handle("PUT /user/password", auth(handlers.UpdatePassword))
	mux.Handle("PATCH /user/password", auth(handlers.UpdatePassword))
	mux.Handle("PUT /user/email", auth(handlers.UpdateEmail))
	mux.Handle("PATCH /user/email", auth(handlers.UpdateEmail))
	mux.Handle("DELETE /user/delete", auth(handlers.DeleteUser))

	// ===== Стеки =====
	mux.Handle("POST /user/stack", auth(handlers.AddUserStack))
	mux.Handle("DELETE /user/stack", auth(handlers.DeleteStack))
	mux.Handle("PUT /user/stack/update", auth(handlers.UpdateStack))
	mux.Handle("PATCH /user/stack/update", auth(handlers.UpdateStack))

	// ===== О себе =====
	mux.Handle("POST /users/about", auth(handlers.DescriptionUser))

	// ===== Проекты =====
	mux.Handle("POST /project", auth(projectHAND.CreatProject))
	mux.HandleFunc("GET /project/public", projectHAND.GetProjectPublik)
	mux.HandleFunc("GET /project/{id}", projectHAND.GetProject)
	mux.Handle("PUT /project/{id}/privacy", auth(projectHAND.UpdatePrivateProject))
	mux.Handle("PATCH /project/{id}/privacy", auth(projectHAND.UpdatePrivateProject))
	mux.Handle("PUT /project/{id}/status", auth(projectHAND.UpdateStatusProject))
	mux.Handle("PATCH /project/{id}/status", auth(projectHAND.UpdateStatusProject))
	mux.Handle("PUT /project/{id}", auth(projectHAND.UpdateProject))
	mux.Handle("DELETE /project/{id}", auth(projectHAND.DeleteProject))

	// ===== Слоты =====
	mux.Handle("PUT /slot/{id}/status", auth(slotHAND.UpdateStatusSlot))
	mux.Handle("PATCH /slot/{id}/status", auth(slotHAND.UpdateStatusSlot))

	// ===== Заявки на слоты =====
	mux.Handle("POST /slot/{id}/apply", auth(slotHAND.ApplicationsSlot))
	mux.Handle("GET /slot/{id}/applications", auth(slotHAND.GetSlotApplications))
	mux.Handle("PUT /slot/{id}/approve/{userID}", auth(slotHAND.ApproveApplication))
	mux.Handle("PATCH /slot/{id}/approve/{userID}", auth(slotHAND.ApproveApplication))
	mux.Handle("PUT /slot/{id}/reject/{userID}", auth(slotHAND.RejectApplication))
	mux.Handle("PATCH /slot/{id}/reject/{userID}", auth(slotHAND.RejectApplication))
	mux.Handle("DELETE /slot/{id}/member", auth(slotHAND.RemoveSlotMember))

	// ===== Запуск =====
	fmt.Println("Сервер запущен: 8080")
	if err := http.ListenAndServe(":8080", service.CORSMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}

//изучить контекст в БД запросах (завершение запроса при зависании)

// ЗАДАЧИ
// 1. написать 2 функции для уникальности логина и никнейма
// 2. Прописать логику сохранения стека в БД + будем брать из БД и отображать на главной странице

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
