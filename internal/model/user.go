package model

type User struct {
	ID        int64  `json:"id"`
	Account   string `json:"account"`
	Username  string `json:"username"`
	Balance   string `json:"balance"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
