package repositories

import (
	"HOTA/internal/models"
	"database/sql"
	"fmt"
)

// перезаписываем стек пользователя: удаляем старый и добавляем новый
// вызываем только внутри транзакции (tx), чтобы при ошибке всё откатилось
func replaceUserStacks(tx *sql.Tx, userID int, stackIDs []int) error {

	// удаляем все старые стеки пользователя
	_, err := tx.Exec(`DELETE FROM user_stacks WHERE user_id = ?`, userID)
	if err != nil {
		return fmt.Errorf("ошибка удаления стеков пользователя: %w", err)
	}

	// добавляем новые стеки по одному
	for _, stackID := range stackIDs {
		_, err := tx.Exec(`INSERT INTO user_stacks (user_id, stack_id) VALUES (?, ?)`, userID, stackID)
		if err != nil {
			return fmt.Errorf("ошибка добавления стека пользователю: %w", err)
		}
	}

	return nil
}

// заменяем стек пользователя целиком (ручка PUT /user/stack/update)
// если что-то пошло не так — у пользователя остаётся старый стек
func ReplaceUserStacks(userID int, stackIDs []int) error {

	// начинаем транзакцию
	tx, err := models.UserDB.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback() // если не дошли до Commit — всё отменится

	err = replaceUserStacks(tx, userID, stackIDs)
	if err != nil {
		return err
	}

	// сохраняем изменения
	return tx.Commit()
}

// проверяем, что все id стеков есть в таблице stacks и не повторяются
func ValidateStacksExist(stackIDs []int) error {

	// сюда складываем id, которые уже проверили, чтобы поймать повтор
	proverennye := map[int]bool{}

	for _, id := range stackIDs {

		// такой id уже был в списке
		if proverennye[id] {
			return fmt.Errorf("стек с id %d указан дважды", id)
		}
		proverennye[id] = true

		// есть ли такой стек в БД
		var exists bool
		query := `SELECT EXISTS(SELECT 1 FROM stacks WHERE id = ?)`
		err := models.UserDB.QueryRow(query, id).Scan(&exists)
		if err != nil {
			return fmt.Errorf("ошибка проверки стека: %w", err)
		}
		if !exists {
			return fmt.Errorf("стек с id %d не найден", id)
		}
	}

	return nil
}
