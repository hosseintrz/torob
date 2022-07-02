package persistence

import (
	"github.com/hosseintrz/torob/user/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	AddUser(user *model.User) (primitive.ObjectID, error)
	GetUser(username string) (*model.User, error)
	DeleteUser(username string) ([]byte, error)
	UpdateUser(user *model.User) (int64, error)
}
