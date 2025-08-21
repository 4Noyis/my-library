package models

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Book    *Book  `json:"book"`
}
