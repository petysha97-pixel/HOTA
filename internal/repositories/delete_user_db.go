package repositories

import (
	"HOTA/internal/models"
	"fmt"
)

// удаляем пользователя вместе со всем, что на него ссылается
// всё делаем в одной транзакции: если хоть один шаг упал — ничего не удалится
func DeleteUser(id int) error {

	// начинаем транзакцию
	tx, err := models.UserDB.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback() // если не дошли до Commit — всё отменится

	// 1. заявки на слоты в проектах пользователя
	_, err = tx.Exec(`DELETE FROM applicationsSlot WHERE slot_id IN (
		SELECT slots.id FROM slots
		INNER JOIN projects ON slots.project_id = projects.id
		WHERE projects.owner_id = ?)`, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления заявок на слоты проектов: %w", err)
	}

	// 2. стеки слотов в проектах пользователя
	_, err = tx.Exec(`DELETE FROM slot_stacks WHERE slot_id IN (
		SELECT slots.id FROM slots
		INNER JOIN projects ON slots.project_id = projects.id
		WHERE projects.owner_id = ?)`, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления стеков слотов проектов: %w", err)
	}

	// 3. слоты в проектах пользователя
	_, err = tx.Exec(`DELETE FROM slots WHERE project_id IN (SELECT id FROM projects WHERE owner_id = ?)`, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления слотов проектов: %w", err)
	}

	// 4. сами проекты пользователя
	_, err = tx.Exec(`DELETE FROM projects WHERE owner_id = ?`, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления проектов: %w", err)
	}

	// 5. заявки пользователя на чужие слоты
	_, err = tx.Exec(`DELETE FROM applicationsSlot WHERE user_id = ?`, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления заявок пользователя: %w", err)
	}

	// 6. чужие слоты, где он исполнитель и работа не сдана, снова открываем
	_, err = tx.Exec(`UPDATE slots SET status = 'open', user_id = NULL WHERE user_id = ? AND status != 'done'`, id)
	if err != nil {
		return fmt.Errorf("ошибка открытия слотов пользователя: %w", err)
	}

	// 7. в сданных слотах (done) просто убираем его, статус не трогаем
	_, err = tx.Exec(`UPDATE slots SET user_id = NULL WHERE user_id = ?`, id)
	if err != nil {
		return fmt.Errorf("ошибка снятия пользователя со сданных слотов: %w", err)
	}

	// 8. стек пользователя
	_, err = tx.Exec(`DELETE FROM user_stacks WHERE user_id = ?`, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления стеков пользователя: %w", err)
	}

	// 9. сам пользователь
	res, err := tx.Exec(`DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении пользователя: %w", err)
	}

	// если удалено 0 строк — такого пользователя не было
	count, _ := res.RowsAffected()
	if count == 0 {
		return fmt.Errorf("пользователь %d не найден", id)
	}

	// сохраняем изменения
	return tx.Commit()
}
