package models

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Book    *Book  `json:"book"`
}

type UserResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
