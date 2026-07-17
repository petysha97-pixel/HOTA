package models

import (
	
)





type RegistrationResponseError struct {
	StatusGlobal string `json:"statusglobal"`
	Error        error `json:"error,omitempty"`
	ID           int    `json:"id,omitempty"`
}