package user

import "errors"

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("O nome é obrigatório")
	}
	if u.Email == "" {
		return errors.New("o email é obrigatório")
	}
	return nil
}
