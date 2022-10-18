package app

import "time"

type Spec struct {
	Name        string    `json:"name,omitempty" bson:"name,omitempty"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"`
	Version     string    `json:"version,omitempty" bson:"version,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
