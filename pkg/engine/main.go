package engine

import (
	"errors"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/updatecli/updateserver/pkg/dashboard"
	"github.com/updatecli/updateserver/pkg/database"
	"github.com/updatecli/updateserver/pkg/server"
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

func (e *Engine) StartServer() {
	server.Run()
}

func (e *Engine) StartRunner() {
	var dashboards []dashboard.Dashboard
	var err error

	if err := database.Connect(e.Options.Database); err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}

	if len(e.Options.Dashboards) > 0 {
		for i := range e.Options.Dashboards {
			d := e.Options.Dashboards[i]
			if err := d.Init(); err != nil {
				logrus.Errorln(err)
			}
			err := d.Save()
			if err != nil {
				logrus.Errorln(err)
			}
		}
	}

	for {
		// Load all dashboard to update them
		if dashboards, err = dashboard.SearchAll(); err != nil {
			logrus.Println(err)
			continue
		}

		// Update Dashboard
		for _, dashboard := range dashboards {

			if err := dashboard.Run(); err != nil {
				logrus.Errorln(err)
				continue
			}
			if err := dashboard.Save(); err != nil {
				logrus.Errorln(err)
				continue
			}
		}

		logrus.Infof("work done, doing a 10 secondes break")
		time.Sleep(10 * time.Second)

	}

}

func (e *Engine) Start() {
	go e.StartRunner()
	e.StartServer()
}
