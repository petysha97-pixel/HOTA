package repositories

import (
	"HOTA/internal/models"
	"strings"
)

var users []models.User
var nextID = 1

// CreateUser Добавляем нового пользователя в слайс пользователей
func CreateUser(user models.User) {
	user.ID = nextID
	users = append(users, user)
	nextID++

}

// FindUserByNickname Фукция поиска по нику
func FindUserByNickname(nickname string) *models.User {

	for _, user := range users {
		if nickname == user.Nickname {
			return &user
		}

	}
	return nil
}

// FindByStack Поиск по стеку
func FindByStack(stack string) []models.User {
	result := make([]models.User, 0)

	for _, user := range users {
		for _, steckITM := range user.Staсk {
			if strings.EqualFold(steckITM, stack) {
				result = append(result, user)
				break
			}
		}

	}
	return result
}

// FindUserByID Поиск по ID
func FindUserByID(id int) *models.User {

	for _, user := range users {
		if id == user.ID {
			return &user
		}
	}
	return nil
}

// Обновление пользователя
func UpdateUser(user models.User) bool {

	for i := range users {
		if users[i].ID == user.ID {
			users[i] = user
			return true
		}
	}
	return false
}

// вывод всех пользователей
func GetAllUsers() []models.User {
	return users
}

// удаление пользователя по ID
func DeleteUser(id int) bool {
	for i := range users {
		if users[i].ID == id {
			users = append(users[:i], users[i+1:]...)
			return true
		}
	}
	return false
}

// статистика (количество пользователей всего)
func CountUser() int {
	return len(users)
}

// статистика (количество пользователей по ролям)
func CountByRole() map[string]int {
	rezult := make(map[string]int) //Backend: 2  Frontend: 1

	for _, user := range users {
		rezult[user.Role]++

	}

	return rezult
}
