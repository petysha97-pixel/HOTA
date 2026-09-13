package service

import (
	"HOTA/internal/models"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Обновление пароля
func ValidateUpdatePass(id int, oldPass, newPass string) error {

	if oldPass == "" {
		return errors.New("текущий пароль не может быть пустым")
	}
	if newPass == "" {
		return errors.New("новый пароль не может быть пустым")
	}

	//Старый пароль из БД
	var storedHash string
	query := "SELECT Password FROM users WHERE id = ?"
	err := models.UserDB.QueryRow(query, id).Scan(&storedHash)
	if err != nil {
		return fmt.Errorf("ошибка получения пароля из БД: %w", err)
	}

	//Хеширование введенного старого пароля
	if err := CheckPassword(storedHash, oldPass); 
	err != nil {
		return errors.New("введенный пароль не верный")
	}

	//Валидация и Хеширование нового пароля
	if err := validatePassword(newPass); err != nil {
		return err
	}

	hashNewPass, err := Hash_password(newPass)
	if err != nil {
		return fmt.Errorf("ошибка хеширования нового пароля: %w", err)
	}

	//Сравнение нового пароля со старым
	if hashNewPass == storedHash {                      //проверить на птоснный false, тк при хешировании соль разная
		return errors.New("новый пароль не может совпадать с текущим")
	}

	//Обновление пароля в БД
	updateQuery := "UPDATE users SET Password = ? WHERE id = ?"
	_, err = models.UserDB.Exec(updateQuery, hashNewPass, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления пароля: %w", err)
	}

	return nil
}

func ValidateNewPasswordOnly(newPass string) error {
	data := models.UpdatePasswors{
		NewPassword: newPass,
	}

	return validation.ValidateStruct(&data, validation.Field(&data.NewPassword, validation.Required, validation.Length(8, 30), validation.By(ValidatePasswors)))
}

// кастомка
func validatePassword(password string) error {

	if !allowedCharsPattern.MatchString(password) {
		return errors.New("пароль должен содержать только буквы, цифры и спецсимволы (!@#$%^&*()-+=), длина 8-30 символов")
	}

	if !hasUpperPattern.MatchString(password) {
		return errors.New("пароль должен содержать хотя бы одну заглавную букву")
	}

	if !hasLowerPattern.MatchString(password) {
		return errors.New("пароль должен содержать хотя бы одну строчную букву")
	}

	if !hasDigitPattern.MatchString(password) {
		return errors.New("пароль должен содержать хотя бы одну цифру")
	}

	if !hasSpecialPattern.MatchString(password) {
		return errors.New("пароль должен содержать хотя бы один спецсимвол (!@#$%^&*()-+=)")
	}

	return nil
}
