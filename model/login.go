package model

type Login struct {
	Username string `json:"userName"`
	Password string `json:"password"`
}

type TokenResp struct {
	AccessToken string `json:"access_token"`
}
