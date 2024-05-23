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
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UserUpdate struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type UserChangePassword struct {
	NewPassword   string `json:"new_password"`
	ReNewPassword string `json:"re_new_password"`
	OldPassword   string `json:"old_password"`
}
