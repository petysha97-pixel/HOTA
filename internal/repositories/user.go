package repositories

import (
	"NOTA/internal/models"
)

var Users []models.User
var nextID = 1

// Добавляем нового пользователя в слайс пользователей
func AppendUser(user models.User) {
	user.ID = nextID
	Users = append(Users, user)
	nextID++
	//Надо сделать переменную nextid = 1 и инкрементировать ее в AppendUser
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

// // Поиск по стеку
// func FindByStack(stack string) []User {
// 	result := make([]User, 0)

// 	for _, user := range users {
// 		for _, steckITM := range user.Staсk {
// 			if strings.EqualFold(steckITM, stack) {
// 				result = append(result, user)
// 			}
// 		}

// 	}
// 	return result
// }

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


