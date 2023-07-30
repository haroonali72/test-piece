package models

import "errors"

// User represents the user data model
type User struct {
	Username   string   `json:"username" validate:"required"`
	ExpiryDate int64    `json:"expiry_date" validate:"required"`
	Outputs    []string `json:"outputs" validate:"required"`
	Password   string   `json:"password" validate:"required"`
}

// Validate checks if the required fields are present in the user model
func (user *User) Validate() error {
	if user.Username == "" {
		return errors.New("username is required")
	}

	if user.ExpiryDate == 0 {
		return errors.New("expiry_date is required")
	}

	if len(user.Outputs) == 0 {
		return errors.New("outputs are required")
	}

	if user.Password == "" {
		return errors.New("password is required")
	}

	return nil
}
