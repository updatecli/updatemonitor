package dashboard

import (
	"context"
	"errors"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/updatecli/updateserver/pkg/app"
	"github.com/updatecli/updateserver/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrFailedUpdatingDashboard error = errors.New("failed updating dashboard")
)

type Project struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bjson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bjson:"name,omitempty"`
	Apps []app.App          `json:"apps,omitempty" bjson:"apps,omitempty"`
}

type Dashboard struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bjson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bjson:"name,omitempty"`
	Projects []Project          `json:"projects,omitempty" bjson:"projects,omitempty"`
}

// loadDashboard query a database to retrieve projects
func Search() ([]Dashboard, error) {

	collection := database.Client.Database(app.DatabaseName).Collection(app.DatabaseCollection)

	var dashboards []Dashboard
	//var searchLimit int64 = 10

	//filter := bson.D{{Key: "id", Value: "xxx"}}
	//findOptions := options.FindOptions{
	//	Limit: &searchLimit,
	//}

	//if err := collection.Find(context.TODO(), filter, findOptions).Decode(&apps); err != nil {
	//	return err
	//}
	cursor, err := collection.Find(context.TODO(), bson.M{})

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

		//var emptyDashboard Dashboard

		//if dashboard == emptyDashboard {
		//	return nil, nil
		//}

		dashboards = append(dashboards, dashboard)
	}

	return dashboards, nil
}

func (d *Dashboard) Run() error {

	for _, project := range d.Projects {
		for _, app := range project.Apps {
			if err := app.Run(); err != nil {
				logrus.Errorln(err)
				continue
			}

			if err := app.Save(); err != nil {
				logrus.Errorln(err)
				continue
			}
		}
	}

	return nil
}
