package app

import (
	"fmt"
	"time"

	updatecliSource "github.com/updatecli/updatecli/pkg/core/pipeline/source"
	"gopkg.in/yaml.v3"
)

const (
	// ErrData is the default error message if version mismatch
	ErrData string = "ERROR"
)

type Data struct {
	Data       Spec   `json:"data,omitempty" bson:"data,omitempty"`
	DataSource string `json:"datasource,omitempty" bson:"datasource,omitempty"`
}

func (d *Data) Run() error {

	currentTime := time.Now().UTC()

	if d.Data.CreatedAt.IsZero() {
		d.Data.CreatedAt = currentTime
		d.Data.UpdatedAt = currentTime
	}

	sourceSpec := updatecliSource.Config{}
	if err := yaml.Unmarshal([]byte(d.DataSource), &sourceSpec); err != nil {
		return err
	}

	if err := sourceSpec.Validate(); err != nil {
		return err
	}

	source := updatecliSource.Source{
		Config: sourceSpec,
	}

	// Disable for testing
	//if d.Data.UpdatedAt.After(currentTime.Add(-10 * time.Minute)) {
	//	logrus.Infof("Data already updated within the latest 10 minutes")
	//	return nil
	//}

	if err := source.Run(); err != nil {
		d.Data.Version = ErrData
		d.Data.UpdatedAt = currentTime
		return fmt.Errorf("failed execute source: %s", err)
	}

	if d.Data.Version != source.Output {
		d.Data.UpdatedAt = currentTime
		d.Data.Version = source.Output
	}

	return nil
}

func (d Data) IsZero() bool {
	var zero Data
	return d == zero
}
