package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Привет, Мир")
}

func main() {
	http.HandleFunc("/", hello)

	fmt.Println("Сервер запущен на :8080")

	http.ListenAndServe(":8080", nil)
}

// https://t.me/GO_Blc
// https://t.me/GO_Blc
// https://t.me/GO_Blc
// goto
// faltruni

//Debian ubunta - на них крутиться сервер
// Obsidian
// func printInfo(user User) {

// 	fmt.Printf("Покупатель %s. Телефон: %s. Адрес: г. %s, %s.\n", user.Name, user.Phone, user.Address.City, user.Address.Street)

// Python + Django 2 - 3 тыс запросов, а Go - 100 тыс запросов
//Tailwind CSS и будстрап
