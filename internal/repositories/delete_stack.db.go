package repositories

import (
	"HOTA/internal/models"
	"fmt"
)

func DeleteStackUser(userID, stackID int) error{

	
	qwery := `DELETE FROM users_stack WHERE user_id = ? AND stack_id = ?`
	_, err := models.UserDB.Exec(qwery,userID, stackID)
	if err != nil {
		return fmt.Errorf("Ошибка удлания стека %w", err)

	}

   return nil
	

}