package user

import (
	"context"
	"fmt"
)

type Service interface {
	Create(ctx context.Context, input CreateUserInput) (UserOutput, error)
	GetAll(ctx context.Context) ([]UserOutput, error)
	GetByID(ctx context.Context, id int) (UserOutput, error)
	Update(ctx context.Context, id int, input UpdateUserInput) (UserOutput, error)
	Delete(ctx context.Context, id int) error
}

type userService struct {
	repo Repository
}

func NewUserService(r Repository) Service {
	return &userService{repo: r}
}

func (s *userService) Create(ctx context.Context, input CreateUserInput) (UserOutput, error) {
	u := User{
		Name:  input.Name,
		Email: input.Email,
	}

	if err := u.Validate(); err != nil {
		return UserOutput{}, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	createdUser, err := s.repo.Create(ctx, u)
	if err != nil {
		return UserOutput{}, err
	}

	return ToUserOutput(createdUser), nil
}

func (s *userService) GetAll(ctx context.Context) ([]UserOutput, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return ToUserOutputList(users), nil
}

func (s *userService) GetByID(ctx context.Context, id int) (UserOutput, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return UserOutput{}, err
	}

	return ToUserOutput(u), nil
}

func (s *userService) Update(ctx context.Context, id int, input UpdateUserInput) (UserOutput, error) {
	req := UpdateUserRequest{
		Name:  input.Name,
		Email: input.Email,
	}

	updatedUser, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return UserOutput{}, err
	}

	return ToUserOutput(updatedUser), nil
}

func (s *userService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
