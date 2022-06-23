package model

import "time"

type User struct {
	ID          string
	FullName    string
	Username    string
	Email       string
	Password    string
	CreatedDate time.Time
	Role        int32
}

func NewUser(id, fullName, username, email, password string, createdDate time.Time, role int32) *User {
	return &User{
		ID:          id,
		FullName:    fullName,
		Username:    username,
		Email:       email,
		Password:    password,
		CreatedDate: createdDate,
		Role:        role,
	}
}
