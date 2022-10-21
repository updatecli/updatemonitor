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
	Spec           Spec   `json:"data,omitempty" bson:"data,omitempty"`
	UpdateManifest string `json:"updatemanifest,omitempty" bson:"updatemanifest,omitempty"`
}

func (d Data) IsZero() bool {
	var zero Data
	return d == zero
}

// Refactor in Updatecli codebase
func (d *Data) RunUpdatePipeline() error {

	if d.UpdateManifest == "" {
		logrus.Infof("no update manifest provided")
		return nil
	}

	currentTime := time.Now().UTC()

	if d.Spec.CreatedAt.IsZero() {
		d.Spec.CreatedAt = currentTime
		d.Spec.UpdatedAt = currentTime
	}

	if d.Spec.UpdatedAt.After(currentTime.Add(-30 * time.Second)) {
		logrus.Debugf("Data updated less than 30 seconds ago, skipping")
		return nil
	}

	pipelineSpec := updatecliConfig.Spec{}
	// To implement, templating
	// Templating is needed to retrieve environment variables
	t := updatecliConfig.Template{}

	templatedUpdateManifest, err := t.New([]byte(d.UpdateManifest))
	if err != nil {
		return err
	}
	//

	if err := yaml.Unmarshal(templatedUpdateManifest, &pipelineSpec); err != nil {
		logrus.Errorf("failed parsing Update manifest - %s", err.Error())
		return err
	}

	pipelineConfig := updatecliConfig.Config{
		Spec: pipelineSpec,
	}

	if err := pipelineConfig.EnsureLocalScm(); err != nil {
		logrus.Errorf("failed generate local scm handler - %s", err.Error())
		return err
	}

	if err := pipelineConfig.Validate(); err != nil {
		logrus.Errorf("failed validating update manifest - %s", err.Error())
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
			if _, err := scm.Handler.Clone(); err != nil {
				logrus.Errorln(err)
			}
		}
	default:
		logrus.Errorf("%d scm configuration detected in update manifest, only one will be considered")
	}

	switch len(pipeline.Sources) {
	case 0:
		// nothing to do
	case 1:
		for i := range pipeline.Sources {
			source := pipeline.Sources[i]
			if err := source.Run(); err != nil {
				d.Spec.Version = ErrData
				d.Spec.UpdatedAt = currentTime
				return fmt.Errorf("failed executing source: %s", err.Error())
			}

			if d.Spec.Version != source.Output {
				d.Spec.Version = source.Output
			}
			d.Spec.UpdatedAt = currentTime

			logrus.Infof("Version %q retrieved at %q", d.Spec.Version, d.Spec.UpdatedAt.String())
		}
	default:
		logrus.Errorf("%d source configuration detected in update manifest, only one will be considered")

	}

	return nil
}
