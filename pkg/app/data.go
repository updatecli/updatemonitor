package app

import (
	"fmt"
	"time"

	updatecliSource "github.com/updatecli/updatecli/pkg/core/pipeline/source"
)

const (
	// ErrData is the default error message if version mismatch
	ErrData string = "ERROR"
)

type Data struct {
	Data       Spec                   `json:"data,omitempty"`
	DataSource updatecliSource.Config `json:"datasource,omitempty"`
}

func (d *Data) Run() error {
	currentTime := time.Now().UTC()

	if len(d.Data.CreateAt) == 0 {
		d.Data.CreateAt = currentTime.String()
	}

	if err := d.DataSource.Validate(); err != nil {
		return err
	}

	source := updatecliSource.Source{
		Config: d.DataSource,
	}

	// ToDO: To retrieve source data from source Database and only updating
	// if newer than 1min

	if err := source.Run(); err != nil {
		d.Data.Version = ErrData
		d.Data.UpdatedAt = currentTime.String()
		return fmt.Errorf("failed execute source: %s", err)
	}

	fmt.Printf("\n\n")

	if d.Data.Version != source.Output {
		d.Data.UpdatedAt = currentTime.String()
		d.Data.Version = source.Output
	}

	return nil
}
