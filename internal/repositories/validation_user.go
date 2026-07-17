package repositories

import (
	"HOTA/internal/models"
	"database/sql"
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var User models.User

// валидация полей при регистрации
func ValidateStruct(User models.User) error {
	return validation.ValidateStruct(&User,

		//Логин (мыло @ + домен)
		validation.Field(&User.Login, validation.Required, is.Email, validation.Length(1, 30), validation.By(UNIK_login)),

		//Пароль (P-p 1-9 1(!-%))
		validation.Field(&User.Password, validation.Required, validation.Length(8, 30), validation.By(ValidatePasswors)),

		// Никнейм (от 2 символов)
		validation.Field(&User.Nickname, validation.Required, validation.Length(2, 20), validation.By(UNIK_Nickname)),

		// Роль обязательна и должна быть одной из строго заданных на фронтенде
		validation.Field(&User.Rolle, validation.Required, validation.In(
			"Frontend", "Backend", "Fullstack", "DevOps")),

		// Стек обязателен. Мы проверяем каждый элемент массива (каждую строку технологии)
		validation.Field(&User.Stack, validation.Required, validation.Length(1, 6), validation.Each(
			validation.Required,
			validation.Length(1, 10), // Каждая технология от 1 до 10 символов
		)),
	)
}

// Только для проверки допустимых символов
var allowedCharsPattern = regexp.MustCompile(`^[A-Za-z0-9!@#$%^&*()\-+=]{8,30}$`)

// Отдельные паттерны для каждого требования
var hasUpperPattern = regexp.MustCompile(`[A-Z]`)
var hasLowerPattern = regexp.MustCompile(`[a-z]`)
var hasDigitPattern = regexp.MustCompile(`[0-9]`)
var hasSpecialPattern = regexp.MustCompile(`[!@#$%^&*()\-+=]`)

func ValidatePasswors(password any) error {

	pass, ok := password.(string)
	if !ok {
		return errors.New("Пароль должен быть строкой")
	}

	if !allowedCharsPattern.MatchString(pass) {
		return errors.New("пароль должен содержать: спецсимвол (!@#$%^&*()-+=)")
	}

	if !hasUpperPattern.MatchString(pass) {
		return errors.New("пароль должен содержать: заглавнй буквы")
	}
	if !hasLowerPattern.MatchString(pass) {
		return errors.New("пароль должен содержать: строчные буквы")
	}
	if !hasDigitPattern.MatchString(pass) {
		return errors.New("пароль должен содержать: хотя бы 1 цифру ")
	}
	if !hasSpecialPattern.MatchString(pass) {
		return errors.New("пароль должен содержать: спецсимвол (!@#$%^&*()-+=)")
	}

	return nil
}

// Уникальность Никнейма
func UNIK_Nickname(value any) error {
	nickname := value.(string)

	var cont int

	err := models.UserDB.QueryRow(
		"SELECT 1 FROM users WHERE Nickname = ? LIMIT 1",
		nickname,
	).Scan(&cont)

	if errors.Is(err, sql.ErrNoRows) {
		// Ник свободен
		return nil
	}

	if err != nil {
		// Настоящая ошибка БД
		return err
	}

	// Ник найден
	return errors.New("ник уже занят")
}

// Уникальность Логина
func UNIK_login(login any) error {

	log, ok := login.(string)
	if !ok {
		return errors.New("Логин должен быть строкой")
	}

	var count int

	query := `SELECT COUNT(1) FROM users WHERE Login = ? LIMIT 1`

	err := models.UserDB.QueryRow(query, log).Scan(&count)

	if errors.Is(err, sql.ErrNoRows) {
		// логин свободен
		return nil
	}

	if err != nil {
		return err //ошиба из БД
	}

	if count > 0 {
		return errors.New("ник уже занят")
	}

	return nil
}
