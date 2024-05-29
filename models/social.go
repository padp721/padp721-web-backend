package models

import "github.com/google/uuid"

type Social struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Url      string    `json:"url"`
	Color    string    `json:"color"`
	IconType string    `json:"icon_type"`
	Icon     string    `json:"icon"`
}

type SocialInput struct {
	Name     string `json:"name" validate:"required,min=1"`
	Url      string `json:"url" validate:"required,min=1"`
	Color    string `json:"color" validate:"required,min=1"`
	IconType string `json:"icon_type" validate:"required,min=1"`
	Icon     string `json:"icon" validate:"required,min=1"`
}
