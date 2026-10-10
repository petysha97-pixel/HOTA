package repositories

import (
	"HOTA/internal/models"
	"fmt"
)

// AppendUser добавляет нового пользователя вместе со стеком одной транзакцией и возвращает его же с ID
func AppendUser(user models.User) (models.User, error) {

	tx, err := models.UserDB.Begin()
	if err != nil {
		return models.User{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() // откатится, если не закоммитим

	userSQL := `
	INSERT INTO users (
		Email,
		Password,
	    Nickname,
		Name,
		Rolle,
		Grade
	   )
	VALUES (?, ?, ?, ?, ?, ?);
	`

	row, err := tx.Exec(
		userSQL,
		user.Email,
		user.Password,
		user.Nickname,
		user.Name,
		user.Rolle,
		user.Grade,
	)
	if err != nil {
		return models.User{}, fmt.Errorf("ошибка сохранения пользователя: %w", err)
	}

	//достаем id пользователя из БД
	id, err := row.LastInsertId()
	if err != nil {
		return models.User{}, fmt.Errorf("ошибка получения id пользователя: %w", err)
	}
	user.ID = int(id)

	//сохраненния связей многие ко многим
	if err := replaceUserStacks(tx, user.ID, user.StackID); err != nil {
		return models.User{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.User{}, fmt.Errorf("commit: %w", err)
	}

	return user, nil
}
