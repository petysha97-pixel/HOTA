package repositories

import (
	"HOTA/internal/models"
	"strings"
)

var Users []models.User
var nextID = 1

// Добавляем нового пользователя в слайс пользователей
func AppendUser(user models.User) models.User{
	user.ID = nextID
	Users = append(Users, user)
	nextID++


	return user
}

// Фукция поиска по нику
func SersheNickname(nickname string) *models.User {

	for _, user := range Users {
		if nickname == user.Nickname {
			return &user
		}

	}
	return nil
}

// Поиск по стеку
func FindByStack(stack string) []models.User {
	result := make([]models.User, 0)

	for _, user := range Users {
		for _, steckITM := range user.Stack {
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

	for _, user := range Users {
		if id == user.ID {
			return &user
		}
	}
	return nil
}

// Обновление пользователя
func UpdateUser(user models.User) bool {

	for i := range Users {
		if Users[i].ID == user.ID {
			Users[i] = user
			return true
		}
	}
	return false
}

// вывод всех пользователей
func GetAllUsers() []models.User {
	return Users
}

// удаление пользователя по ID
func DeleteUser(id int) bool {
	for i := range Users {
		if Users[i].ID == id {
			Users = append(Users[:i], Users[i+1:]...)
			return true
		}
	}
	return false
}

// статистика (количество пользователей всего)
func CountUser() int {
	return len(Users)
}

// статистика (количество пользователей по ролям)
func CountByRole() map[string]int {
	rezult := make(map[string]int) //Backend: 2  Frontend: 1

	for _, user := range Users {
		rezult[user.Role]++

	}

	return rezult
}
