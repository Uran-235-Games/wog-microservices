package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func (m *MongoDB) Connect(uri string) error {
	docs := "https://www.mongodb.com/docs/drivers/go/current/"

	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"See: " + docs +
			"usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	m.Client = client
	m.Database = client.Database("wogmain")

	return nil
}

func (m *MongoDB) Disconnect() error {
	return m.Client.Disconnect(context.TODO())
}
