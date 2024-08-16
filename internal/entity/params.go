package entity

// Response Модель ответа сервера
type Response struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Errors  string `json:"errors"`
}

type Env struct {
	AppPort int
	DBUrl   string
}

type Filter struct {
	Cond       []map[string]string `json:"cond"`
	Conditions map[string][]string `json:"conditions"`
	Limit      int                 `json:"limit"`
}
