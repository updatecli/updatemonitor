package dashboard

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/updatecli/updateserver/pkg/app"
	"github.com/updatecli/updateserver/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DatabaseName       string = "dashboard"
	DatabaseCollection string = "dashboard"
)

var (
	ErrFailedUpdatingDashboard error = errors.New("failed updating dashboard")
)

type Project struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Apps []app.App          `json:"apps,omitempty" bson:"apps,omitempty"`
}

type Dashboard struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Owner    string             `json:"owner,omitempty" bson:"owner,omitempty"`
	Projects []Project          `json:"projects,omitempty" bson:"projects,omitempty"`
}

// loadDashboard query a database to retrieve projects
func SearchAll() ([]Dashboard, error) {

	collection := database.Client.Database(DatabaseName).Collection(DatabaseCollection)

	var dashboards []Dashboard
	var searchLimit int64 = 10

	filter := bson.D{}
	findOptions := options.FindOptions{
		Limit: &searchLimit,
		Sort: bson.D{
			{
				Key:   "updatedAt",
				Value: -1,
			},
		},
	}

	cursor, err := collection.Find(context.TODO(), filter, &findOptions)

	if err != nil {
		return nil, err
	}

	ctx := context.TODO()
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var dashboard Dashboard
		if err = cursor.Decode(&dashboard); err != nil {
			log.Fatal(err)
			break
		}

		dashboards = append(dashboards, dashboard)
	}

	return dashboards, nil
}

func DeletebyID(ID string) (*mongo.DeleteResult, error) {

	collection := database.Client.Database(DatabaseName).Collection(DatabaseCollection)

	var dashboard Dashboard

	objId, err := primitive.ObjectIDFromHex(ID)
	//
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.DeleteOne(ctx, bson.M{"id": objId})
	defer cancel()

	if err != nil {
		logrus.Println(dashboard)
		logrus.Errorln(err)
		return nil, err
	}

	return result, nil
}

func SearchbyID(ID string) (Dashboard, error) {

	collection := database.Client.Database(DatabaseName).Collection(DatabaseCollection)

	var dashboard Dashboard

	objId, err := primitive.ObjectIDFromHex(ID)
	//
	if err != nil {
		return dashboard, err
	}

	filter := bson.D{{Key: "_id", Value: objId}}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	result := collection.FindOne(ctx, filter)
	defer cancel()

	err = result.Decode(&dashboard)
	if err != nil {
		logrus.Println(dashboard)
		logrus.Errorln(err)
		return dashboard, err
	}

	return dashboard, nil
}

func (d *Dashboard) Init() error {

	if d.ID.IsZero() {
		d.ID = primitive.NewObjectID()
	}

	if d.Owner == "" {
		d.Owner = "default"
	}

	for i := range d.Projects {
		project := d.Projects[i]
		for j := range project.Apps {
			app := project.Apps[j]
			err := app.Init()
			if err != nil {
				return err
			}
			project.Apps[j] = app
		}
		d.Projects[i] = project
	}

	return nil

}

func (d *Dashboard) Run() error {

	for i, project := range d.Projects {
		for j, app := range project.Apps {
			if err := app.Run(); err != nil {
				logrus.Errorln(err)
				continue
			}

			if err := app.Save(); err != nil {
				logrus.Errorln(err)
				continue
			}
			project.Apps[j] = app
		}

		d.Projects[i] = project
	}

	return nil
}

func (d Dashboard) Save() error {

	collection := database.Client.Database(DatabaseName).Collection(DatabaseCollection)

	// https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/write-operations/upsert/
	///***
	// */

	if d.ID.IsZero() {
		return fmt.Errorf("dashboard ID is not defined")
	}

	filter := bson.D{
		{
			Key:   "name",
			Value: d.Name,
		},
		{
			Key:   "owner",
			Value: d.Owner,
		}}
	update := bson.D{
		{
			Key: "$set", Value: bson.D{
				{
					Key: "projects", Value: d.Projects,
				},
				{
					Key: "owner", Value: d.Owner,
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

	fmt.Printf("Number of documents updated: %d\n", result.ModifiedCount)
	fmt.Printf("Number of documents upserted: %d\n", result.UpsertedCount)
	return nil
}

func (d Dashboard) SaveByID(ID string) (*mongo.UpdateResult, error) {

	collection := database.Client.Database(DatabaseName).Collection(DatabaseCollection)

	// https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/write-operations/upsert/
	///***
	// */

	objId, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{
			Key:   "_id",
			Value: objId,
		},
	}
	update := bson.D{
		{
			Key: "$set", Value: bson.D{
				{
					Key: "projects", Value: d.Projects,
				},
				{
					Key: "owner", Value: d.Owner,
				},
			},
		},
	}
	opts := options.Update().SetUpsert(true)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	defer cancel()

	if err != nil {
		return nil, err
	}

	fmt.Printf("Number of documents updated: %d\n", result.ModifiedCount)
	fmt.Printf("Number of documents upserted: %d\n", result.UpsertedCount)
	return result, nil
}
