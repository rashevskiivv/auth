package entity

// Token cross layers token entity.
type Token struct {
	ID     *int64 `json:"id,omitempty"`
	Token  string `json:"token"`
	UserID int64  `json:"user_id"`
}
