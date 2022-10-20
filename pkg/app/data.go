package app

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	updatecliConfig "github.com/updatecli/updatecli/pkg/core/config"
	updatecliPipeline "github.com/updatecli/updatecli/pkg/core/pipeline"
	"gopkg.in/yaml.v3"
)

const (
	// ErrData is the default error message if version mismatch
	ErrData string = "ERROR"
)

type Data struct {
	Data           Spec   `json:"data,omitempty" bson:"data,omitempty"`
	UpdateManifest string `json:"updatemanifest,omitempty" bson:"updatemanifest,omitempty"`
}

func (d Data) IsZero() bool {
	var zero Data
	return d == zero
}

func (d *Data) RunUpdatePipeline() error {

	if d.UpdateManifest == "" {
		logrus.Infof("no update manifest provided")
		return nil
	}

	currentTime := time.Now().UTC()

	if d.Data.CreatedAt.IsZero() {
		d.Data.CreatedAt = currentTime
		d.Data.UpdatedAt = currentTime
	}

	pipelineSpec := updatecliConfig.Spec{}
	// To implement, templating

	//

	if err := yaml.Unmarshal([]byte(d.UpdateManifest), &pipelineSpec); err != nil {
		logrus.Errorln("failed parsing Update manifest - %s", err.Error())
		return err
	}

	pipelineConfig := updatecliConfig.Config{
		Spec: pipelineSpec,
	}

	if err := pipelineConfig.EnsureLocalScm(); err != nil {
		logrus.Errorln("failed generate local scm handler - %s", err.Error())
		return err
	}

	if err := pipelineConfig.Validate(); err != nil {
		logrus.Errorln("failed validating update manifest - %s", err.Error())
		return err
	}

	pipeline := updatecliPipeline.Pipeline{}
	if err := pipeline.Init(&pipelineConfig, updatecliPipeline.Options{}); err != nil {
		logrus.Errorf("failed initiating Update pipeline - %s", err.Error())
	}

	// InitSCM
	// We should check that we only have one scm defined

	switch len(pipeline.SCMs) {
	case 0:
		// skip nothing to do
	case 1:
		for i := range pipeline.SCMs {
			scm := pipeline.SCMs[i]
			scm.Handler.Clone()
		}
	default:
		logrus.Errorf("%d scm configuration detected in update manifest, only one will be considered")
	}

	switch len(pipeline.SCMs) {
	case 0:
		// nothing to do
	case 1:
		for i := range pipeline.Sources {
			source := pipeline.Sources[i]
			if err := source.Run(); err != nil {
				d.Data.Version = ErrData
				d.Data.UpdatedAt = currentTime
				return fmt.Errorf("failed executing source: %s", err.Error())
			}

			if d.Data.Version != source.Output {
				d.Data.UpdatedAt = currentTime
				d.Data.Version = source.Output
			}

			logrus.Infof("Version %q retrieved at %q", d.Data.Version, d.Data.UpdatedAt.String())
		}
	default:
		logrus.Errorf("%d source configuration detected in update manifest, only one will be considered")

	}

	return nil
}
