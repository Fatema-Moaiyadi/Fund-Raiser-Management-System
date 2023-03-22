package models

type UserLoginRequest struct {
	EmailID  string `json:"email_id"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Code int `json:"code"`
	Data struct {
		AuthToken string `json:"auth_token"`
	} `json:"data"`
}
