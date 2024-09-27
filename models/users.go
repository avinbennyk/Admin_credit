package models

import (
	"errors"
)

// User model struct
type User struct {
	Email   string `json:"email" gorm:"primaryKey"`
	Credits int    `json:"credits" gorm:"not null;default:0"`
	Paused  bool   `json:"paused" gorm:"not null;default:false"`
	Role    string `json:"role" gorm:"not null;default:'user'"` // 'admin' or 'user'
}

// IncrementCredits adds a specified amount to the user's credits
func (u *User) IncrementCredits(amount int) {
	u.Credits += amount
}

// DecrementCredits subtracts a specified amount from the user's credits
func (u *User) DecrementCredits(amount int) error {
	if u.Credits-amount < 0 {
		return errors.New("not enough credits")
	}
	u.Credits -= amount
	return nil
}

// PauseAccount pauses the user account
func (u *User) PauseAccount() {
	u.Paused = true
}

// UnpauseAccount reactivates the user account
func (u *User) UnpauseAccount() {
	u.Paused = false
}

// IsAdmin checks if the user has an admin role
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}
