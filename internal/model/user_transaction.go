package model

type UserTransaction struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Type      string `json:"type"`
	Amount    string `json:"amount"`
	InOrOut   string `json:"in_or_out"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
