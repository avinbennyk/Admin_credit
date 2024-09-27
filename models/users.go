package models

import (
	"errors"
)

type User struct {
	Email   string `json:"email" gorm:"primaryKey"`
	Credits int    `json:"credits" gorm:"not null;default:0"`
	Paused  bool   `json:"paused" gorm:"not null;default:false"`
	Role    string `json:"role" gorm:"not null;default:'user'"`
}

func (u *User) IncrementCredits(amount int) {
	u.Credits += amount
}

func (u *User) DecrementCredits(amount int) error {
	if u.Credits-amount < 0 {
		return errors.New("not enough credits")
	}
	u.Credits -= amount
	return nil
}

func (u *User) PauseAccount() {
	u.Paused = true
}

func (u *User) UnpauseAccount() {
	u.Paused = false
}

func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}
