package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Store struct {
	ID      primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name    string
	OwnerId string
	Url     string
	City    string
}

func NewStore(name, ownerId, url, city string) *Store {
	return &Store{
		Name:    name,
		OwnerId: ownerId,
		Url:     url,
		City:    city,
	}
}
