package entity

// User cross layers user entity.
type User struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	INN      *string `json:"inn"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
}
