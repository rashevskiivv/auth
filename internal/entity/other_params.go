package entity

import (
	"fmt"
	"strings"
)

type Response struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Errors  string `json:"errors"`
}

type RequestUtils struct {
	ID           string
	Token        string
	WhichRequest string
}

func (r *RequestUtils) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return fmt.Errorf("id must be set and not be empty as a header")
	}
	if strings.TrimSpace(r.Token) == "" {
		return fmt.Errorf("token must be set and not be empty as a header")
	}
	return nil
}
