package database

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	URIEnvVariableName string = "UPDATEMONITOR_DB_URI"
)

var (
	Client *mongo.Client

	URI string = os.Getenv(URIEnvVariableName)
)

type Options struct {
	// URI defines the mongodb URI
	URI string
}

func Connect(o Options) error {

	var err error

	if o.URI != "" {
		if URI != "" {
			logrus.Debugf("URI %q defined from setting file override the value from environment variable  %q", o.URI, URIEnvVariableName)
		}
		URI = o.URI
	}

	clientOptions := options.Client().ApplyURI(URI)

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
