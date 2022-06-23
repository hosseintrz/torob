package persistence

import (
	"context"
	"encoding/json"
	"github.com/hosseintrz/torob/user/internal/db"
	"github.com/hosseintrz/torob/user/internal/model"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	DB    = "torob"
	USERS = "users"
)

type MongoLayer struct {
	db *db.MongoDB
}

func NewMongoLayer(db *db.MongoDB) *MongoLayer {
	return &MongoLayer{db: db}
}

func (m *MongoLayer) AddUser(user *model.User) ([]byte, error) {
	res, err := m.db.Client.Database(DB).Collection(USERS).InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	insertedId, err := json.Marshal(res.InsertedID)
	return insertedId, err
}

func (m *MongoLayer) GetUser(username string) (*model.User, error) {
	res := m.db.Client.Database(DB).Collection(USERS).FindOne(context.Background(), bson.D{{"username", username}})
	if err := res.Err(); err != nil {
		return nil, err
	}
	var user model.User
	err := res.Decode(&user)
	return &user, err
}

func (m *MongoLayer) DeleteUser(username string) ([]byte, error) {
	_, err := m.db.Client.Database(DB).Collection(USERS).DeleteOne(context.Background(), bson.D{{"username", username}})
	if err != nil {
		return nil, err
	}
	return []byte(username), nil
}
func (m *MongoLayer) UpdateUser(user *model.User) (int64, error) {
	res, err := m.db.Client.Database(DB).Collection(USERS).UpdateOne(
		context.Background(),
		bson.D{{"username", user.Username}},
		user)
	if err != nil {
		return 0, err
	}
	cnt := res.ModifiedCount
	return cnt, nil
}
