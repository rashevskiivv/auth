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
