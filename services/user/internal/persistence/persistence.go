package persistence

import "github.com/hosseintrz/torob/user/internal/model"

type UserRepository interface {
	AddUser(user *model.User) ([]byte, error)
	GetUser(username string) (*model.User, error)
	DeleteUser(username string) ([]byte, error)
	UpdateUser(user *model.User) (int64, error)
}
