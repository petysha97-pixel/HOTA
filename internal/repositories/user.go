package repositories

import (
	"HOTA/internal/models"
	"encoding/json"
	"log"
	"os"
	"strings"
)

var Users models.Users
var nextID = 1

// init пытается прочитать существующий файл при запуске приложения
func init() {
	readUserData()
}

// updateUserData обновляет файл с пользователями или создает его
func updateUserData() error {
	Users.Count = len(Users.Data)
	data, err := json.Marshal(&Users)
	if err != nil {
		return err
	}
	err = os.WriteFile("../users_data.json", data, 0666)
	if err != nil {
		return err
	}

	return nil
}

// readUserData читает файл с пользователями
func readUserData() error {
	file, err := os.ReadFile("../users_data.json")
	if err != nil {
		return err
	}
	json.Unmarshal(file, &Users)

	return nil
}

// Добавляем нового пользователя в слайс пользователей
func AppendUser(user models.User) {
	user.ID = Users.Data[len(Users.Data)-1].ID + 1 // ID последнего добавленного пользователя +1
	Users.Data = append(Users.Data, user)
	nextID++
	err := updateUserData()
	if err != nil {
		log.Fatal(err)
	}
}

// Фукция поиска по нику
func SersheNickname(nickname string) *models.User {

	for _, user := range Users.Data {
		if nickname == user.Nickname {
			return &user
		}

	}
	return nil
}

// Поиск по стеку
func FindByStack(stack string) []models.User {
	result := make([]models.User, 0)

	for _, user := range Users.Data {
		for _, steckITM := range user.Staсk {
			if strings.EqualFold(steckITM, stack) {
				result = append(result, user)
				break
			}
		}

	}
	return result
}

// Поиск по ID
func SersheID(id int) *models.User {

	for _, user := range Users.Data {
		if id == user.ID {
			return &user
		}
	}
	return nil
}

// Обновление пользователя
func UpdateUser(user models.User) bool {

	for i := range Users.Data {
		if Users.Data[i].ID == user.ID {
			Users.Data[i] = user
			return true
		}
	}
	err := updateUserData()
	if err != nil {
		log.Fatal(err)
	}
	return false
}

// вывод всех пользователей
func GetAllUsers() []models.User {
	return Users.Data
}

// удаление пользователя по ID
func DeleteUser(id int) bool {
	for i := range Users.Data {
		if Users.Data[i].ID == id {
			Users.Data = append(Users.Data[:i], Users.Data[i+1:]...)
			return true
		}
	}
	Users.Count = len(Users.Data)
	err := updateUserData()
	if err != nil {
		log.Fatal(err)
	}
	return false
}

// статистика (количество пользователей всего)
func CountUser() int {
	return Users.Count
}

// статистика (количество пользователей по ролям)
func CountByRole() map[string]int {
	rezult := make(map[string]int) //Backend: 2  Frontend: 1

	for _, user := range Users.Data {
		rezult[user.Role]++
	}

	return rezult
}