package slot

import (
	"HOTA/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

var ErrApplicationConflict = errors.New("не удалось подать отклик: слот или заявка изменились")

// Создаём отклик, если набор ещё открыт.
func CreateApplication(slotID, userID int, message string) error {
	// SELECT вернёт строку для вставки только при открытом наборе.
	result, err := models.UserDB.Exec(`
		INSERT INTO applicationsSlot (slot_id, user_id, status, message)
		SELECT s.id, ?, 'pending', ?
		FROM slots s
		JOIN projects p ON p.id = s.project_id
		WHERE s.id = ?
		  AND s.status IN ('open', 'review')
		  AND p.status != 'finished'
		ON CONFLICT(slot_id, user_id) DO NOTHING
	`, userID, message, slotID)
	if err != nil {
		return fmt.Errorf("создание заявки: %w", err)
	}

	// Если строка не добавилась, набор закрыт или отклик уже существует.
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("проверка создания заявки: %w", err)
	}
	if count == 0 {
		return ErrApplicationConflict
	}
	return nil
}

// ищем заявку пользователя на слот, которая ещё на рассмотрении
func GetUserCheckingApplications(slotID, userID int) (*models.AplicationSlot, error) {
	query := `SELECT id, slot_id, user_id, status, message, created_at FROM applicationsSlot
	WHERE slot_id = ?
	AND user_id = ?
	AND status = 'pending'`

	var app models.AplicationSlot
	err := models.UserDB.QueryRow(query, slotID, userID).Scan(&app.ID, &app.SlotID, &app.UserID, &app.Status, &app.Message, &app.CreatAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска заявки: %w", err)
	}

	return &app, nil
}

// достаём все заявки по ID слота
func GetApplicationsBySlot(slotID int) ([]models.AplicationSlot, error) {
	query := `SELECT id, slot_id, user_id, status, message, created_at FROM applicationsSlot WHERE slot_id = ?`
	rows, err := models.UserDB.Query(query, slotID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения заявок: %w", err)
	}
	defer rows.Close()

	app := []models.AplicationSlot{}
	for rows.Next() {
		var a models.AplicationSlot
		err := rows.Scan(&a.ID, &a.SlotID, &a.UserID, &a.Status, &a.Message, &a.CreatAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования: %w", err)
		}
		app = append(app, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения заявок: %w", err)
	}
	return app, nil
}

// обновляем статус заявки
func UpdateApplicationStatus(appID int, status string) error {
	query := `UPDATE applicationsSlot SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := models.UserDB.Exec(query, status, appID)
	if err != nil {
		return fmt.Errorf("ошибка обновления заявки: %w", err)
	}
	return nil
}

// отклоняем все остальные заявки на слот, кроме указанной
func RejectOtherApplications(slotID, excludeAppID int) error {
	query := `UPDATE applicationsSlot SET status = 'rejected', updated_at = CURRENT_TIMESTAMP WHERE slot_id = ? AND id != ? AND status = 'pending'`
	_, err := models.UserDB.Exec(query, slotID, excludeAppID)
	if err != nil {
		return fmt.Errorf("ошибка отклонения заявок: %w", err)
	}
	return nil
}

// проверка на то, что в слоте утверждён разработчик
func ApprovedAplecation(slotID int) (bool, error) {
	query := `SELECT COUNT(*) FROM applicationsSlot
	WHERE slot_id = ?
	AND status = 'approved'`

	var count int

	err := models.UserDB.QueryRow(query, slotID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("ошибка проверки утверждённого разработчика на слот: %w", err)
	}

	return count > 0, nil
}

// при переоткрытии снимаем принятую заявку
func ResetApprovedApplication(slotID int) error {
	query := `UPDATE applicationsSlot SET status = 'rejected', updated_at = CURRENT_TIMESTAMP WHERE slot_id = ? AND status = 'approved'`
	_, err := models.UserDB.Exec(query, slotID)
	if err != nil {
		return fmt.Errorf("ошибка снятия заявки: %w", err)
	}
	return nil
}

// снимаем утвержденного разработчика со слота
func RemoveApprovedMember(slotID int) (int64, error) {
	// Начинаем общую операцию из нескольких SQL-запросов.
	tx, err := models.UserDB.Begin()
	if err != nil {
		return 0, fmt.Errorf("начало транзакции: %w", err)
	}
	defer tx.Rollback()

	// Снимаем утверждённую заявку.
	res, err := tx.Exec(`
		UPDATE applicationsSlot
		SET status = 'removed',
		    updated_at = CURRENT_TIMESTAMP
		WHERE slot_id = ?
		  AND status = 'approved'
	`, slotID)
	if err != nil {
		return 0, fmt.Errorf("снятие заявки: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("подсчёт изменённых заявок: %w", err)
	}

	if count == 0 {
		return 0, errors.New("нет утверждённого разработчика для снятия")
	}

	if count != 1 {
		return 0, errors.New("на слоте несколько утверждённых заявок")
	}

	// Убираем исполнителя и открываем слот.
	res, err = tx.Exec(`
		UPDATE slots
		SET user_id = NULL,
		    status = 'open'
		WHERE id = ?
		  AND status IN ('open', 'review')
	`, slotID)
	if err != nil {
		return 0, fmt.Errorf("освобождение слота: %w", err)
	}

	changed, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("подсчёт изменённых слотов: %w", err)
	}

	if changed != 1 {
		return 0, errors.New("слот не найден или его нельзя освободить")
	}

	// Только здесь сохраняются оба изменения.
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("сохранение транзакции: %w", err)
	}

	return count, nil
}

// любая заявка пользователя на слот (в таблице она одна: UNIQUE(slot_id, user_id))
// и сколько минут прошло с последней смены её статуса
func GetUserApplication(slotID, userID int) (*models.AplicationSlot, int, error) {
	query := `SELECT id, slot_id, user_id, status, message, created_at,
	CAST((julianday('now') - julianday(COALESCE(updated_at, created_at))) * 24 * 60 AS INTEGER)
	FROM applicationsSlot
	WHERE slot_id = ? AND user_id = ?`

	var app models.AplicationSlot
	var minutes int
	err := models.UserDB.QueryRow(query, slotID, userID).Scan(&app.ID, &app.SlotID, &app.UserID, &app.Status, &app.Message, &app.CreatAt, &minutes)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка поиска заявки: %w", err)
	}

	return &app, minutes, nil
}

// Возвращаем старый отклик в очередь, если прошла пауза.
func ReopenApplication(appID int, message string, cooldownMinutes int) error {
	// Проверяем состояние заявки и слота в момент обновления.
	result, err := models.UserDB.Exec(`
		UPDATE applicationsSlot
		SET status = 'pending', message = ?,
		    created_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		  AND status IN ('rejected', 'removed')
		  AND (julianday('now') - julianday(COALESCE(updated_at, created_at))) * 1440 >= ?
		  AND EXISTS (
		      SELECT 1
		      FROM slots s
		      JOIN projects p ON p.id = s.project_id
		      WHERE s.id = applicationsSlot.slot_id
		        AND s.status IN ('open', 'review')
		        AND p.status != 'finished'
		  )
	`, message, appID, cooldownMinutes)
	if err != nil {
		return fmt.Errorf("повторная подача заявки: %w", err)
	}

	// Если условие перестало выполняться, отклик остаётся как был.
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("проверка обновления заявки: %w", err)
	}
	if count == 0 {
		return ErrApplicationConflict
	}
	return nil
}

var ErrApprovalConflict = errors.New("слот занят или заявка уже изменилась")

// утверждение кондидата и утверждения на слот
// Остальные заявки отклоняются.
func ApproveSlotApplication(slotID, applicantID, ownerID int) error {
	tx, err := models.UserDB.Begin()
	if err != nil {
		return fmt.Errorf("начало транзакции: %w", err)
	}
	defer tx.Rollback()

	// Назначаем кандидата, только если слот всё ещё свободен,
	// а проект принадлежит этому владельцу и не завершён.
	res, err := tx.Exec(`
	UPDATE slots
	SET user_id = ?,
    status = 'close'
	WHERE id = ?
	  AND (
	      (status = 'open' AND user_id IS NULL)
	      OR
	      (status = 'review' AND user_id = ?)
	  )
	  AND EXISTS (
	      SELECT 1 FROM projects
	      WHERE id = slots.project_id
	        AND owner_id = ?
	        AND status != 'finished'
	  )
`, applicantID, slotID, applicantID, ownerID)
	if err != nil {
		return fmt.Errorf("назначение исполнителя: %w", err)
	}

	changed, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("подсчёт изменённых слотов: %w", err)
	}
	if changed != 1 {
		return ErrApprovalConflict
	}

	// Утверждаем только заявку, которая ещё ждёт решения.
	res, err = tx.Exec(`
		UPDATE applicationsSlot
		SET status = 'approved', updated_at = CURRENT_TIMESTAMP
		WHERE slot_id = ? AND user_id = ? AND status = 'pending'
	`, slotID, applicantID)
	if err != nil {
		return fmt.Errorf("утверждение заявки: %w", err)
	}

	changed, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("подсчёт изменённых заявок: %w", err)
	}
	if changed != 1 {
		return ErrApprovalConflict
	}

	// Отклоняем оставшиеся заявки, которые ещё ждут решения.
	_, err = tx.Exec(`
		UPDATE applicationsSlot
		SET status = 'rejected', updated_at = CURRENT_TIMESTAMP
		WHERE slot_id = ? AND status = 'pending'
	`, slotID)
	if err != nil {
		return fmt.Errorf("отклонение остальных заявок: %w", err)
	}

	return tx.Commit()
}
