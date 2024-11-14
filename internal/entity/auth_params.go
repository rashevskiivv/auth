package entity

type TokenFilter struct {
	ID     []string `json:"id"`
	Token  []string `json:"token"`
	UserID []string `json:"user_id"`
	Limit  int32    `json:"limit"`
}

type RegisterInput struct {
	User
}

type RegisterOutput struct {
	Token
}

type AuthenticateInput struct {
	User
}

type AuthenticateOutput struct {
	Token
}
