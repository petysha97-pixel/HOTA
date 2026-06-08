package repositories

import (
	"HOTA/src/internal/models"
	"strings"
)

var Users []models.User
var nextID = 1

// ========== SETTERS ==========

func Add(user models.User) {
	user.Id = nextID
	Users = append(Users, user)
	nextID++
}

func Update(user models.User) bool {
	for i := range Users {
		if Users[i].Id == user.Id {
			Users[i] = user
			return true
		}
	}

	return false
}

func Delete(id int) bool {
	for i := range Users {
		if Users[i].Id == id {
			Users = append(Users[:i], Users[i+1:]...)
			return true
		}
	}

	return false
}

// ========== GETTERS ==========

func GetById(id int) *models.User {
	for _, user := range Users {
		if id == user.Id {
			return &user
		}
	}

	return nil
}

func GetByNickname(nickname string) *models.User {
	for _, user := range Users {
		if nickname == user.Nickname {
			return &user
		}
	}

	return nil
}

func GetByStack(stack string) []models.User {
	result := make([]models.User, 0)

	for _, user := range Users {
		for _, stackItem := range user.Stack {
			if strings.EqualFold(stackItem, stack) {
				result = append(result, user)
				break
			}
		}
	}

	return result
}

func CountByRole() map[string]int {
	result := make(map[string]int)

	for _, user := range Users {
		result[user.Role]++
	}

	return result
}
