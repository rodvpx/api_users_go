package user

type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserInput struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

type UserOutput struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ToUserOutput(u User) UserOutput {
	return UserOutput{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func ToUserOutputList(users []User) []UserOutput {
	outputs := make([]UserOutput, len(users))
	for i, u := range users {
		outputs[i] = ToUserOutput(u)
	}
	return outputs
}
