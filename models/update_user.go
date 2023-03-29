package models

type UpdateUser struct {
	FirstName string `db:"first_name,omitempty" json:"first_name,omitempty"`
	LastName  string `db:"last_name,omitempty" json:"last_name,omitempty"`
	Password  string `db:"password,omitempty" json:"password,omitempty"`
}

type UpdateUserResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		UpdatedInfo UpdateUser `json:"updated_info"`
	} `json:"data"`
}
