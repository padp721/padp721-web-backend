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
	Name     string `json:"name"`
	Url      string `json:"url"`
	Color    string `json:"color"`
	IconType string `json:"icon_type"`
	Icon     string `json:"icon"`
}
