package models

type Repository struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	CreateAt    string `json:"create_at"`
	UpdateAt    string `json:"update_at"`
}

type RepositoryCreate struct {
	UserId      int    `json:"user_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Path        string `json:"path" validate:"required"`
}

type RepositoryUpdate struct {
	Id          int    `json:"id" validate:"required"`
	UserId      int    `json:"user_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Path        string `json:"path" validate:"required"`
}
