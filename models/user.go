package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UserID    int
	ChatID    int64
	FirstName string
	LastName  string
	Username  string
}

func (*User) TableName() string {
	return "user"
}

func (u *User) IsEqual(user *User) bool {
	if u.UserID != user.UserID || u.ChatID != user.ChatID ||
		u.FirstName != user.FirstName || u.LastName != user.LastName ||
		u.Username != user.Username {
		return false
	}
	return true
}
