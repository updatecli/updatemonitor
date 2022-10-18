package app

import (
	"github.com/google/uuid"
	updatecliSource "github.com/updatecli/updatecli/pkg/core/pipeline/source"
)

type Source struct {
	ID        uuid.UUID              `json:"_id,omitempty" bson:"_id,omitempty"`
	Result    string                 `json:"result,omitempty" bson:"result,omitempty"`
	CreateAt  string                 `json:"createAt,omitempty" bson:"createAt,omitempty"`
	UpdatedAt string                 `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	Config    updatecliSource.Config `json:"datasource,omitempty" bson:"datasource,omitempty"`
	Checksum  string                 `json:"checksum,omitempty" bson:"checksum,omitempty"`
	Counter   int                    `json:"counter,omitempty" bson:"counter,omitempty"`
}

func (s Source) Save() error {
	return nil
}

func (s *Source) Load() error {
	return nil
}
