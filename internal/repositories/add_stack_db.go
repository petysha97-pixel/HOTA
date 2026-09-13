package repositories

import (
	"HOTA/internal/models"
	"fmt"
)

func AddStackUser(userID, stackID int) error{

	//проверяем есть ли стек такой в БД
	var ex bool

	qwery := `SELECT EXISTS(SELECT 1 FROM stacks WHERE id = ?)`
	err := models.UserDB.QueryRow(qwery, stackID).Scan(&ex)
	if err != nil {
		return fmt.Errorf("ошибка проверки стека %w", err)
	}

	if ex == false {
		return fmt.Errorf("Стек не найден в БД %w", err)
	}

	//Добавляем связи между юзером и стеком (если она ex=труе)
	qwery = `INSERT OR IGNORE INTO user_stacks (user_id, stack_id) 	VALUES (?, ?)`
	_, err = models.UserDB.Exec(qwery, userID, stackID)

	if err != nil {
		return fmt.Errorf("ошибка добавления стека %w", err)
	}

	return nil

}
