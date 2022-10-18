package database

import (
	"context"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = Client.Ping(ctx, nil)
	defer cancel()

	if err != nil {
		logrus.Errorf("database ping test failed: %s", err)
		return err
	}

	logrus.Infoln("database connected")

	return nil
}
