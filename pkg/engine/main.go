package engine

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
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

		//

		logrus.Infof("work done, doing a 5 minutes break")
		time.Sleep(5 * time.Minute)

	}
}
