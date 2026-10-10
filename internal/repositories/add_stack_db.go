package repositories

import (
	"HOTA/internal/models"
	"errors"
	"fmt"
)

// ошибки стека пользователя, хендлер по ним выбирает ответ
var ErrStackNotFound = errors.New("стек не найден")
var ErrStackAlreadyAdded = errors.New("эта технология уже есть в стеке")
var ErrStackLimit = errors.New("в стеке уже 6 технологий")
var ErrStackLast = errors.New("в стеке должна остаться хотя бы одна технология")

// сколько технологий в стеке у пользователя
func CountUserStacks(userID int) (int, error) {
	var count int
	err := models.UserDB.QueryRow(`SELECT COUNT(*) FROM user_stacks WHERE user_id = ?`, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка подсчёта стеков пользователя: %w", err)
	}
	return count, nil
}

// добавляем одну технологию в стек пользователя (не больше 6)
func AddStackUser(userID, stackID int, description string) error {

	//проверяем есть ли стек такой в БД
	var ex bool
	qwery := `SELECT EXISTS(SELECT 1 FROM stacks WHERE id = ?)`
	err := models.UserDB.QueryRow(qwery, stackID).Scan(&ex)
	if err != nil {
		return fmt.Errorf("ошибка проверки стека %w", err)
	}
	if !ex {
		return ErrStackNotFound
	}

	//может, эта технология уже есть у пользователя
	qwery = `SELECT EXISTS(SELECT 1 FROM user_stacks WHERE user_id = ? AND stack_id = ?)`
	err = models.UserDB.QueryRow(qwery, userID, stackID).Scan(&ex)
	if err != nil {
		return fmt.Errorf("ошибка проверки стека пользователя %w", err)
	}
	if ex {
		return ErrStackAlreadyAdded
	}

	//не больше 6 технологий
	count, err := CountUserStacks(userID)
	if err != nil {
		return err
	}
	if count >= 6 {
		return ErrStackLimit
	}

	//Добавляем связь между юзером и стеком
	qwery = `INSERT INTO user_stacks (user_id, stack_id, description) VALUES (?, ?, ?)`
	_, err = models.UserDB.Exec(qwery, userID, stackID, description)
	if err != nil {
		return fmt.Errorf("ошибка добавления стека %w", err)
	}

	return nil
}

// меняем описание опыта у одной технологии пользователя
func UpdateStackInfo(userID, stackID int, description string) error {

	query := `UPDATE user_stacks SET description = ? WHERE user_id = ? AND stack_id = ?`
	result, err := models.UserDB.Exec(query, description, userID, stackID)
	if err != nil {
		return fmt.Errorf("ошибка изменения стека пользователя: %w", err)
	}

	// 0 строк — такой технологии нет в стеке пользователя
	count, _ := result.RowsAffected()
	if count == 0 {
		return ErrStackNotFound
	}

	return nil
}
