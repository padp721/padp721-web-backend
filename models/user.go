package models

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
}

type UserRegister struct {
	Username string `json:"username" validate:"required,min=1"`
	Name     string `json:"name" validate:"required,min=1"`
	Email    string `json:"email" validate:"required,min=1"`
	Phone    string `json:"phone" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=1"`
}

type UserUpdate struct {
	Username string `json:"username" validate:"required,min=1"`
	Name     string `json:"name" validate:"required,min=1"`
	Email    string `json:"email" validate:"required,min=1"`
	Phone    string `json:"phone" validate:"required,min=1"`
}

type UserChangePassword struct {
	NewPassword   string `json:"new_password" validate:"required,min=1"`
	ReNewPassword string `json:"re_new_password" validate:"required,min=1"`
	OldPassword   string `json:"old_password" validate:"required,min=1"`
}
