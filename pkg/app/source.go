package app

import (
	"github.com/google/uuid"
	updatecliSource "github.com/updatecli/updatecli/pkg/core/pipeline/source"
)

type Source struct {
	ID        uuid.UUID
	Result    string
	CreateAt  string                 `json:"createAt,omitempty"`
	UpdatedAt string                 `json:"updatedAt,omitempty"`
	Config    updatecliSource.Config `json:"datasource,omitempty"`
	Checksum  string                 `json:"checksum,omitempty"`
	Counter   int                    `json:"counter,omitempty"`
}

func (s Source) Save() error {
	return nil
}

func (s *Source) Load() error {
	return nil
}
