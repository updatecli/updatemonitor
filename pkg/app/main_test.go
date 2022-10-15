package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/updatecli/updatecli/pkg/core/pipeline/resource"
	updatecliSource "github.com/updatecli/updatecli/pkg/core/pipeline/source"
	"github.com/updatecli/updatecli/pkg/plugins/resources/dockerdigest"
)

func TestRun(t *testing.T) {
	testdata := []struct {
		name            string
		app             App
		expectedVersion string
	}{
		{
			name: "Basic test",
			app: App{
				Expected: Data{
					Data: Spec{
						Title:       "Expected",
						Description: "This is a description",
					},
					DataSource: updatecliSource.Config{
						ResourceConfig: resource.ResourceConfig{
							Kind: "dockerdigest",
							Spec: dockerdigest.Spec{
								Image:        "updatecli/updatecli",
								Tag:          "v0.33.0",
								Architecture: "amd64",
							},
						},
					},
				},
				Current: Data{
					Data: Spec{
						Title:       "Current",
						Description: "This is a description",
					},
					DataSource: updatecliSource.Config{
						ResourceConfig: resource.ResourceConfig{
							Kind: "dockerdigest",
							Spec: dockerdigest.Spec{
								Image:        "updatecli/updatecli",
								Tag:          "v0.33.0",
								Architecture: "amd64",
							},
						},
					},
				},
			},
			expectedVersion: "sha256:cd2225127073f3a778f33d31f798c3014375b04725214a940b9f000fa41d8339",
		},
	}

	for _, tt := range testdata {

		t.Run(tt.name, func(t *testing.T) {

			err := tt.app.Run()

			require.NoError(t, err)

			assert.Equal(t, tt.expectedVersion, tt.app.Expected.Data.Version)

		})
	}

}
