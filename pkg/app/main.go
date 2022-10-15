package app

type Spec struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
	CreateAt    string `json:"createAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
}
