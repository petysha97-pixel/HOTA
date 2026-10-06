package models

type RegistrationResponseError struct {
	StatusGlobal string `json:"statusGlobal"`
	Error        error  `json:"error,omitempty"`
	ID           int    `json:"id,omitempty"`
}
