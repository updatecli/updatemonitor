package app

//import (
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//)

//func TestRun(t *testing.T) {
//	testdata := []struct {
//		name            string
//		expected        Data
//		current         Data
//		expectedVersion string
//	}{
//		{
//			name: "Basic test",
//			expected: Data{
//				Data: Spec{
//					Name:        "Expected",
//					Description: "This is a description",
//				},
//				DataSource: `kind: dockerdigest
//spec:
//    image: "updatecli/updatecli"
//    tag: "v0.33.0"
//    architecture: "amd64"
//`,
//			},
//			current: Data{
//				Data: Spec{
//					Name:        "Current",
//					Description: "This is a description",
//				},
//				DataSource: `kind: dockerdigest
//spec:
//    image: "updatecli/updatecli"
//    tag: "v0.33.0"
//    architecture: "amd64"
//`,
//			},
//			expectedVersion: "sha256:cd2225127073f3a778f33d31f798c3014375b04725214a940b9f000fa41d8339",
//		},
//	}
//
//	for _, tt := range testdata {
//
//		t.Run(tt.name, func(t *testing.T) {
//
//			a, err := Init(tt.current, tt.expected)
//
//			require.NoError(t, err)
//
//			err = a.Run()
//			require.NoError(t, err)
//
//			assert.Equal(t, tt.expectedVersion, a.Expected.Data.Version)
//
//		})
//	}
//
//}
//
