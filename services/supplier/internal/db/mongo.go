package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDB struct {
	opts   *options.ClientOptions
	Client *mongo.Client
}

func NewMongoDB(conn string) *MongoDB {
	serverApiOpts := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(conn).SetServerAPIOptions(serverApiOpts)
	return &MongoDB{
		opts: opts,
	}
}

func (m *MongoDB) Connect() error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, m.opts)
	if err != nil {
		return err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("error connecting to mongo server")
		return err
	}
	m.Client = client
	return nil
}
