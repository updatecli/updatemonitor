package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/updatecli/updateserver/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DatabaseName       string = "app"
	DatabaseCollection string = "app"
	SUCCESS            int    = iota // 0
	WARNING                          // 1
	FAILED                           // 2
)

type App struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bjson:"_id,omitempty"`
	Current     Data               `json:"current,omitempty" bjson:"current,omitempty"`
	CreatedAt   time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Expected    Data               `json:"expected,omitempty" bjson:"expected,omitempty"`
	Status      int                `json:"status,omitempty" bjson:"status,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	UpdatedAt   time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func New(inputCurrent, inputExpected Data) (App, error) {
	var a App
	a.ID = primitive.NewObjectID()

	currentTime := time.Now().UTC()
	a.CreatedAt = currentTime
	a.UpdatedAt = currentTime

	a.Expected = inputExpected
	a.Current = inputCurrent
	return a, nil
}

func (a *App) Run() error {
	errs := []error{}

	logrus.Infof("Updating App %q", a.ID.String())

	// Init App ID if not set
	if a.ID.IsZero() {
		a.ID = primitive.NewObjectID()
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

	return nil
}

func SearchApps() ([]App, error) {

	var apps []App

	collection := database.Client.Database(DatabaseName).Collection(DatabaseCollection)

	cursor, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	ctx := context.TODO()
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var app App
		if err = cursor.Decode(&app); err != nil {
			log.Fatal(err)
			break
		}

		var emptyApp App

		if app == emptyApp {
			return nil, nil
		}

		apps = append(apps, App{
			ID:       app.ID,
			Status:   app.Status,
			Expected: app.Expected,
			Current:  app.Current,
		})
	}

	return apps, nil
}

func (a App) Save() error {

	collection := database.Client.Database(DatabaseName).Collection(DatabaseCollection)

	//filter := bson.M{"_id": a.ID}
	//result, err := collection.ReplaceOne(context.TODO(), filter, a)

	result, err := collection.InsertOne(context.TODO(), a)

	if err != nil {
		logrus.Errorf("failed inserting app document in database - %s", err)
		return err
	}

	// https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/write-operations/upsert/
	///***
	// */
	//filter := bson.D{{"type", "Oolong"}}
	//update := bson.D{{"$set", bson.D{{"rating", 8}}}}
	//opts := options.Update().SetUpsert(true)
	//result, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("Number of documents updated: %v\n", result.ModifiedCount)
	//fmt.Printf("Number of documents upserted: %v\n", result.UpsertedCount)

	logrus.Infof("successfully inserted document in database - %s\n", result)
	return nil
}
