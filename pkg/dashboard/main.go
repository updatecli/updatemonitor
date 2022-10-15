package dashboard

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/updatecli/updateserver/pkg/app"
)

var (
	ErrFailedUpdatingDashboard error = errors.New("failed updating dashboard")
)

type Project struct {
	Name string
	Apps []app.App
}

type Dashboard struct {
	Name     string
	Projects []Project
}

// loadDashboard query a database to retrieve projects
func (d *Dashboard) loadDashboard() error {
	return nil
}

// saveDashboard insert result in the database
func (d *Dashboard) saveDashboard() error {
	return nil
}

func (d *Dashboard) updateDashboard() error {
	errs := []error{}

	for _, project := range d.Projects {
		for _, app := range project.Apps {
			if err := app.Run(); err != nil {
				errs = append(errs, err)
				continue
			}
		}
	}

	for _, err := range errs {
		logrus.Errorln(err)
		return ErrFailedUpdatingDashboard
	}

	return nil
}

func (d *Dashboard) Run() error {

	if err := d.loadDashboard(); err != nil {
		return err
	}

	if err := d.updateDashboard(); err != nil {
		return err
	}

	if err := d.saveDashboard(); err != nil {
		return err
	}
	return nil
}
