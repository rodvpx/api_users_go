package user

import "context"

type Repository interface {
	Create(ctx context.Context, u User) (User, error)
	GetAll(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id int) (User, error)
	Update(ctx context.Context, id int, req UpdateUserRequest) (User, error)
	Delete(ctx context.Context, id int) error
}
