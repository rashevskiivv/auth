package entity

type Response struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Errors  string `json:"errors"`
}

type Env struct {
	AppPort int
	DBUrl   string
}

type UserFilter struct {
	ID    []string `json:"id"`
	Email []string `json:"email"`
	Name  []string `json:"name"`
	Limit int32    `json:"limit"`
}

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
