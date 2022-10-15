package app

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/updatecli/updateserver/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	appDatabaseName       string = "app"
	appDatabaseCollection string = "app"
	SUCCESS               int    = iota // 0
	WARNING                             // 1
	FAILED                              // 2
)

type App struct {
	ID       uuid.UUID `json:"id"`
	Expected Data      `json:"expected,omitempty"`
	Current  Data      `json:"current,omitempty"`
	Status   int       `json:"status,omitempty"`
}

func (a *App) Run() error {
	errs := []error{}

	var emptyUUID uuid.UUID

	// Init App ID if not set
	if a.ID.String() == emptyUUID.String() {
		var space uuid.UUID
		a.ID = uuid.NewHash(
			md5.New(),
			space,
			[]byte(fmt.Sprintf("%v", a)),
			3)
	}

	err := a.Current.Run()
	if err != nil {
		errs = append(errs, fmt.Errorf("current - %s", err))
	}

	err = a.Expected.Run()
	if err != nil {
		errs = append(errs, fmt.Errorf("expected - %s", err))
	}

	switch a.Expected.Data.Version == a.Current.Data.Version {
	case true:
		a.Status = SUCCESS
	case false:
		a.Status = WARNING
	default:
		a.Status = FAILED
	}

	if len(errs) > 0 {
		for _, err := range errs {
			logrus.Errorln(err)
		}
		return fmt.Errorf("failed execute source: %s", err)
	}

	// Save back to database
	if err := a.Save(); err != nil {
		return err
	}

	return nil
}

func SearchMany() ([]App, error) {

	collection := database.Client.Database(appDatabaseName).Collection(appDatabaseCollection)

	var apps []App
	var searchLimit int64 = 10

	filter := bson.D{{Key: "id", Value: "xxx"}}
	findOptions := options.FindOptions{
		Limit: &searchLimit,
	}

	//if err := collection.Find(context.TODO(), filter, findOptions).Decode(&apps); err != nil {
	//	return err
	//}
	cursor, err := collection.Find(context.TODO(), filter, &findOptions)

	if err != nil {
		return nil, err
	}

	cursor.Decode(apps)
	return apps, nil
}

func (a *App) Save() error {

	collection := database.Client.Database(appDatabaseName).Collection(appDatabaseCollection)

	result, err := collection.InsertOne(context.TODO(), a)

	if err != nil {
		logrus.Errorf("failed inserting app document in database - %s", err)
		return err
	}

	logrus.Infof("successfully insert document in database - %s\n", result)
	return nil
}
