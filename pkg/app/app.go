package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/updatecli/updatemonitor/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DatabaseName       string = "app"
	DatabaseCollection string = "app"
	SUCCESS            int    = iota // 0
	WARNING                          // 1
	FAILED                           // 2
)

type App struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Spec        []Data             `json:"spec,omitempty" bson:"spec,omitempty"`
	CreatedAt   time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Status      int                `json:"status,omitempty" bson:"status,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	UpdatedAt   time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func (a *App) Init() error {

	if a.ID.IsZero() {
		a.ID = primitive.NewObjectID()
	}

	currentTime := time.Now().UTC()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = currentTime
		a.UpdatedAt = currentTime
	}

	if a.UpdatedAt.IsZero() {
		a.UpdatedAt = currentTime
	}

	return nil
}

func (a *App) Run() error {
	errs := []error{}

	// Init App ID if not set
	if a.ID.IsZero() {
		a.ID = primitive.NewObjectID()
	}

	if len(a.Spec) == 0 {
		return nil
	}

	logrus.Debugf("Checking App %q\n", a.ID.String())

	for i := range a.Spec {
		d := a.Spec[i]
		err := d.RunUpdatePipeline()
		if err != nil {
			errs = append(errs, fmt.Errorf("current - %s", err))
		}
		a.Spec[i] = d
	}

	matching := true
	foundValue := ""
	for i := range a.Spec {
		if i == 0 {
			foundValue = a.Spec[i].Version
		}
		if foundValue != a.Spec[i].Version {
			logrus.Printf("Value %q mismatch with %q", foundValue, a.Spec[i].Version)
			matching = false
			break
		}
	}

	switch matching {
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
		return fmt.Errorf("failed running app %q - %q", a.Name, a.ID.String())
	}

	return nil
}

func SearchApps() ([]App, error) {
	// https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/sort/

	var apps []App

	collection := database.Client.Database(DatabaseName).Collection(DatabaseCollection)

	filter := bson.M{}
	opts := options.Find().SetSort(
		bson.D{
			{
				Key:   "updatedAt",
				Value: -1},
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, filter, opts)
	defer cancel()

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var app App
		if err = cursor.Decode(&app); err != nil {
			log.Fatal(err)
			break
		}

		if len(app.Spec) == 0 {
			return nil, nil
		}

		apps = append(apps, App{
			ID:     app.ID,
			Status: app.Status,
			Spec:   app.Spec,
		})
	}

	return apps, nil
}

func (a App) Save() error {

	collection := database.Client.Database(DatabaseName).Collection(DatabaseCollection)

	if a.ID.IsZero() {
		return fmt.Errorf("dashboard ID is not defined")
	}

	filter := bson.D{
		{
			Key:   "_id",
			Value: a.ID,
		},
	}
	update := bson.D{
		{
			Key: "$set", Value: bson.D{
				{
					Key: "data", Value: a.Spec,
				},
				{
					Key: "updatedAt", Value: a.UpdatedAt,
				},
			},
		},
	}
	opts := options.Update().SetUpsert(true)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	defer cancel()
	if err != nil {
		return err
	}
	logrus.Debugf("Number of documents updated: %d\n", result.ModifiedCount)
	logrus.Debugf("Number of documents upserted: %d\n", result.UpsertedCount)

	return nil
}
