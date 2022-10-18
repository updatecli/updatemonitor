package engine

import (
	"errors"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/updatecli/updateserver/pkg/app"
	"github.com/updatecli/updateserver/pkg/dashboard"
	"github.com/updatecli/updateserver/pkg/database"
)

var (
	ErrEngineRunFailed error = errors.New("something went wrong in server engine")
)

type Options struct {
	// Engine Options
	Dashboards []dashboard.Dashboard
	Database   database.Options
}

type Engine struct {
	Options Options
}

func (e *Engine) Start() error {

	var dashboards []dashboard.Dashboard
	var err error

	if err := database.Connect(e.Options.Database); err != nil {
		logrus.Errorln(err)
		return err
	}

	if len(e.Options.Dashboards) > 0 {
		dashboards = append(dashboards, e.Options.Dashboards...)
	}

	for {
		for _, dashboard := range dashboards {
			if err := dashboard.Run(); err != nil {
				logrus.Errorln(err)
				continue
			}
		}

		if dashboards, err = dashboard.Search(); err != nil {
			logrus.Println(err)
			os.Exit(1)
		}

		var apps []app.App
		apps, err = app.SearchApps()
		if err != nil {
			break
		}

		for {

			for i := range apps {
				apps[i].Run()
			}
			apps, err = app.SearchApps()
			if err != nil {
				logrus.Errorln(err)
				break
			}

			logrus.Infof("work done, doing a 10 secondes break")
			time.Sleep(10 * time.Second)

		}

		logrus.Infof("work done, doing a 10 secondes break")
		time.Sleep(10 * time.Second)

	}
	return nil
}
