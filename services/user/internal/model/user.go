package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	FullName    string
	Username    string
	Email       string
	Password    string
	CreatedDate time.Time
	Role        int32
}

func NewUser(fullName, username, email, password string, createdDate time.Time, role int32) *User {
	return &User{
		FullName:    fullName,
		Username:    username,
		Email:       email,
		Password:    password,
		CreatedDate: createdDate,
		Role:        role,
	}
}
