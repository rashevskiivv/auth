package entity

import (
	"fmt"
	"strconv"
	"strings"
	"tax-auth/internal/usecase"
)

type UserFilter struct {
	ID    []string `json:"id,omitempty"`
	Email []string `json:"email,omitempty"`
	Name  []string `json:"name,omitempty"`
	Limit int32    `json:"limit,omitempty"`
}

func (f *UserFilter) Validate() error {
	if len(f.ID) > 0 {
		for i, s := range f.ID {
			if len(s) == 0 {
				return fmt.Errorf("%v. id can not be empty", i)
			}
			_, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("%v. id is not integer", i) //todo export errors and use errors.Is to return correct codes
			}
		}
	}

	if len(f.Email) > 0 {
		for i, e := range f.Email {
			if len(e) == 0 {
				return fmt.Errorf("%v. email can not be empty", i)
			}
			emailOk, err := usecase.ValidateEmail(e)
			if err != nil {
				return err
			}
			if !emailOk {
				return fmt.Errorf("%v. email is not email", i)
			}
		}
	}

	if len(f.Name) > 0 {
		for i, n := range f.Name {
			if len(n) == 0 {
				return fmt.Errorf("%v. name can not be empty", i)
			}
			dropFound := strings.Contains(n, "drop")
			if dropFound {
				return fmt.Errorf("%v. name contains \"drop\". It is restricted", i)
			}
			deleteFound := strings.Contains(n, "delete")
			if deleteFound {
				return fmt.Errorf("%v. name contains \"delete\". It is restricted", i)
			}
		}
	}

	if f.Limit < 0 {
		return fmt.Errorf("limit can not be negative")
	}

	return nil
}

type GetUsersInput struct {
	Filter UserFilter `json:"omitempty"`
}

type GetUsersOutput struct {
	Response []User `json:"users"`
}

type UpdateUsersInput struct {
	Model  User       `json:"user"`
	Filter UserFilter `json:"filter"`
}

type DeleteUsersInput struct {
	Filter UserFilter
}
