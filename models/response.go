package models

type Response struct {
	Message string `json:"message"`
}

type ResponseData struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
