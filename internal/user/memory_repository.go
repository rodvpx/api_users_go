package user

import (
	"context"
)

type inMemoryRepository struct {
	users []User
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		users: []User{},
	}
}

func (m *inMemoryRepository) Create(ctx context.Context, u User) (User, error) {
	for _, existingUser := range m.users {
		if existingUser.Email == u.Email {
			return User{}, ErrEmailAlreadyExists
		}
	}

	u.ID = len(m.users) + 1

	m.users = append(m.users, u)
	return u, nil
}

func (m *inMemoryRepository) GetAll(ctx context.Context) ([]User, error) {
	return m.users, nil
}

func (m *inMemoryRepository) GetByID(ctx context.Context, id int) (User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return User{}, ErrUserNotFound
}

func (m *inMemoryRepository) Update(ctx context.Context, id int, req UpdateUserRequest) (User, error) {
	for i, u := range m.users {
		if u.ID == id {
			if req.Name != nil {
				m.users[i].Name = *req.Name
			}
			if req.Email != nil {
				m.users[i].Email = *req.Email
			}

			return m.users[i], nil
		}
	}

	return User{}, ErrUserNotFound
}

func (m *inMemoryRepository) Delete(ctx context.Context, id int) error {

	for i, u := range m.users {
		if u.ID == id {
			m.users = append(m.users[:i], m.users[i+1:]...)

			return nil
		}
	}

	return ErrUserNotFound
}
