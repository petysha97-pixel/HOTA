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

	//регистрация
	mux.Handle("/user", service.POSTMiddleware(http.HandlerFunc(handlers.NewUser)))
	mux.Handle("/stack", service.POSTMiddleware(http.HandlerFunc(handlers.GetStacks)))

	//Авторизация
	mux.Handle("/user/auth", service.POSTMiddleware(http.HandlerFunc(handlers.Auth)))

	//Главный профиль
	mux.Handle("/profile", service.POSTMiddleware(service.JWTMiddleware(http.HandlerFunc((handlers.GetUser)))))

	//CRUD
	mux.Handle("/user/update", service.PutPatchMiddleware(service.JWTMiddleware(http.HandlerFunc(handlers.UpdateUser))))
	mux.Handle("/user/password", service.PutPatchMiddleware(service.JWTMiddleware(http.HandlerFunc(handlers.UpdatePassword))))
	mux.Handle("/user/email", service.PutPatchMiddleware(service.JWTMiddleware(http.HandlerFunc(handlers.UpdateEmail))))
	mux.Handle("/user/delete", service.DeleteMiddleware(service.JWTMiddleware(http.HandlerFunc(handlers.DeleteUser))))

	//стеки
	mux.Handle("/user/stack", service.POSTMiddleware(service.JWTMiddleware(http.HandlerFunc(handlers.AddUserStack))))
	mux.Handle("/user/stack/update", service.PutPatchMiddleware(service.JWTMiddleware(http.HandlerFunc(handlers.UpdateStack))))

	//"О себе"
	mux.Handle("/users/about", service.POSTMiddleware(service.JWTMiddleware(http.HandlerFunc(handlers.DescriptionUser))))

	//ПРОСМОТР ПРОФИЛЯ ДРУГОГО РАЗРАБОТЧИКА
	mux.Handle("/profile/other/{id}", service.GETMiddleware(http.HandlerFunc((handlers.GetUserOtherProfile))))

	//"Умный поиск" - доработать
	mux.Handle("/users/searche", service.GETMiddleware(http.HandlerFunc(handlers.SearcheUsers)))

	//создание проектов со слотами
	//проекты (в проекты создается хотя бы один слот обязательно (слот создается вместе с проектом)))
	mux.Handle("/project", service.POSTMiddleware(service.JWTMiddleware(http.HandlerFunc(projectHAND.CreatProject))))
	// получение проекта по ID
	mux.Handle("/project/{id}", service.GETMiddleware(http.HandlerFunc(projectHAND.GetProject)))
	//Поулчаем все публичные проекты
	mux.HandleFunc("/project/publiс", projectHAND.GetProjectPubliс)
	//смена статуса проекта
	mux.Handle("/project/{id}/status", service.PutPatchMiddleware(service.JWTMiddleware(http.HandlerFunc(projectHAND.UpdateStatusProject))))

	//смена статуса слота
	mux.Handle("/slot/{id}/status", service.PutPatchMiddleware(service.JWTMiddleware(http.HandlerFunc(slotHAND.UpdateStatusSlot))))

	//заявки на слоты

	//подача заявки
	mux.Handle("/slot/{id}/apply", service.POSTMiddleware(service.JWTMiddleware(http.HandlerFunc(slotHAND.ApplicationsSlot))))
	//просмотр заявок на слот
	mux.Handle("/slot/{id}/applications", service.GETMiddleware(service.JWTMiddleware(http.HandlerFunc(slotHAND.GetSlotApplications))))
	//пирнять заявку
	mux.Handle("/slot/{id}/approve/{userID}", service.PutPatchMiddleware(service.JWTMiddleware(http.HandlerFunc(slotHAND.ApproveApplication))))
	//отклонить заявку
	mux.Handle("/slot/{id}/reject/{userID}", service.PutPatchMiddleware(service.JWTMiddleware(http.HandlerFunc(slotHAND.RejectApplication))))

	// 3. Оборачиваем весь роутер в наше CORS Middleware
	fmt.Println("Сервер запущен: 8080")
	http.ListenAndServe(":8080", service.CORSMiddleware(mux))

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
