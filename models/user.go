package models

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type UserSecret struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
