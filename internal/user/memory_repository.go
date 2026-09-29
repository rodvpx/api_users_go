package user

import "errors"

type inMemoryRepository struct {
	users []User
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		users: []User{},
	}
}

func (m *inMemoryRepository) Create(u User) error {
	for _, existingUser := range m.users {
		if existingUser.Email == u.Email {
			return errors.New("Este email já está em uso")
		}
	}

	m.users = append(m.users, u)
	return nil
}

func (m *inMemoryRepository) GetAll() ([]User, error) {
	return m.users, nil
}

func (m *inMemoryRepository) GetByID(id int) (User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return User{}, errors.New("Usuário não encontrado")
}

func (m *inMemoryRepository) Update(id int, newData User) (User, error) {

	for i, u := range m.users {
		if u.ID == id {
			m.users[i].Name = newData.Name
			m.users[i].Email = newData.Email

			return m.users[i], nil
		}
	}

	return User{}, errors.New("Usuário não encontrado")
}

func (m *inMemoryRepository) Delete(id int) error {

	for i, u := range m.users {
		if u.ID == id {
			m.users = append(m.users[:i], m.users[i+1:]...)

			return nil
		}
	}

	return errors.New("Usuário nao encontrado")
}
