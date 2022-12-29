package model

type Authorization struct {
	Status   string `json:"status"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Token    string `json:"token"`
}
