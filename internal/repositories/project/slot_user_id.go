package project

import (
	"database/sql"
)

// исполнитель слота в БД может быть пустым (NULL), если на слот ещё никого не назначили
// sql.NullInt64 умеет хранить такое пустое значение: Valid = false значит NULL
// превращаем его в *int: nil — исполнителя нет, иначе — его id
func UserIDOrNil(userID sql.NullInt64) *int {
	if !userID.Valid {
		return nil
	}
	id := int(userID.Int64)
	return &id
}
