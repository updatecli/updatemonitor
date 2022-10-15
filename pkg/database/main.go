package database

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
)

type Options struct {
	// URI defines the mongodb URI
	URI string
}

func Connect(o Options) error {

	var err error

	clientOptions := options.Client().ApplyURI(o.URI)

	Client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		logrus.Errorf("database connection failed: %s", err)
		return err
	}

	err = Client.Ping(context.TODO(), nil)
	if err != nil {
		logrus.Errorf("database ping test failed: %s", err)
		return err
	}

	logrus.Infof("database connected")

	return nil
}
