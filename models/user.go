package models

import "time"

const UserTable string = "Users"

type User struct {
	UserID      string
	MailAddress string
	FirstName   string
	LastName    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (v *User) Table() string {
	return UserTable
}
