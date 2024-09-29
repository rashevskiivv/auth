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
	Limit int      `json:"limit"`
}

type RegisterInput struct {
	User
}

type AuthenticateInput struct {
	Password string
	Hash     string
}

type RegisterOutput struct {
	Token
}

type AuthenticateOutput struct {
	Token
}
