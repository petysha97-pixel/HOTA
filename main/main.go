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
	err := godotenv.Load("../.env")
	if err != nil {
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

	// метод указывается прямо в шаблоне пути (Go 1.22+):
	// на другой метод роутер сам ответит 405, отдельные method-middleware не нужны
	auth := func(h http.HandlerFunc) http.Handler { return service.JWTMiddleware(h) }

	//регистрация
	mux.HandleFunc("POST /user", handlers.NewUser)
	mux.HandleFunc("GET /stack", handlers.GetStacks)

	//Авторизация
	mux.HandleFunc("POST /user/auth", handlers.Auth)

	//Главный профиль
	mux.Handle("/profile", service.GETMiddleware(service.JWTMiddleware(http.HandlerFunc((handlers.GetUser)))))

	//CRUD
	mux.Handle("PUT /user/update", auth(handlers.UpdateUser))
	mux.Handle("PATCH /user/update", auth(handlers.UpdateUser))
	mux.Handle("PUT /user/password", auth(handlers.UpdatePassword))
	mux.Handle("PATCH /user/password", auth(handlers.UpdatePassword))
	mux.Handle("PUT /user/email", auth(handlers.UpdateEmail))
	mux.Handle("PATCH /user/email", auth(handlers.UpdateEmail))
	mux.Handle("DELETE /user/delete", auth(handlers.DeleteUser))

	//стеки
	mux.Handle("POST /user/stack", auth(handlers.AddUserStack))
	mux.Handle("DELETE /user/stack", auth(handlers.DeleteStack))
	mux.Handle("PUT /user/stack/update", auth(handlers.UpdateStack))
	mux.Handle("PATCH /user/stack/update", auth(handlers.UpdateStack))

	//"О себе"
	mux.Handle("POST /users/about", auth(handlers.DescriptionUser))

	//все разработчики
	mux.HandleFunc("GET /users", handlers.GetAllUser)

	//ПРОСМОТР ПРОФИЛЯ ДРУГОГО РАЗРАБОТЧИКА
	mux.HandleFunc("GET /profile/other/{id}", handlers.GetUserOtherProfile)

	//"Умный поиск" - доработать
	mux.HandleFunc("GET /users/searche", handlers.SearcheUsers)

	//создание проектов со слотами
	//проекты (в проекты создается хотя бы один слот обязательно (слот создается вместе с проектом)))
	mux.Handle("/project", service.POSTMiddleware(service.JWTMiddleware(http.HandlerFunc(projectHAND.CreatProject))))
	// получение проекта по ID
	mux.Handle("/project/{id}", service.GETMiddleware(http.HandlerFunc(projectHAND.GetProject)))
	//Поулчаем все публичные проекты
	mux.HandleFunc("GET /project/public", projectHAND.GetProjectPublik)
	// получение проекта по ID (токен необязателен: нужен только для приватных проектов)
	mux.Handle("GET /project/{id}", service.OptionalJWTMiddleware(http.HandlerFunc(projectHAND.GetProject)))

	//смена статуса проекта
	mux.Handle("PUT /project/{id}/status", auth(projectHAND.UpdateStatusProject))
	mux.Handle("PATCH /project/{id}/status", auth(projectHAND.UpdateStatusProject))
	//смена статуса слота
	mux.Handle("/slot/{id}/status", service.PutPatchMiddleware(service.JWTMiddleware(http.HandlerFunc(slotHAND.UpdateStatusSlot))))
	mux.Handle("PUT /slot/{id}/status", auth(slotHAND.UpdateStatusSlot))
	mux.Handle("PATCH /slot/{id}/status", auth(slotHAND.UpdateStatusSlot))



	//заявки на слоты
	

	//подача заявки
	mux.Handle("POST /slot/{id}/apply", auth(slotHAND.ApplicationsSlot))
	//просмотр заявок на слот

	mux.Handle("GET /slot/{id}/applications", auth(slotHAND.GetSlotApplications))
	//пирнять заявку
	mux.Handle("PUT /slot/{id}/approve/{userID}", auth(slotHAND.ApproveApplication))
	mux.Handle("PATCH /slot/{id}/approve/{userID}", auth(slotHAND.ApproveApplication))
	//отклонить заявку
	mux.Handle("PUT /slot/{id}/reject/{userID}", auth(slotHAND.RejectApplication))
	mux.Handle("PATCH /slot/{id}/reject/{userID}", auth(slotHAND.RejectApplication))

	// 3. Оборачиваем весь роутер в наше CORS Middleware
	fmt.Println("Сервер запущен: 8080")
	http.ListenAndServe(":8080", service.CORSMiddleware(mux))
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
