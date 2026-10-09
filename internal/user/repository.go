package user

type Repository interface {
	Create(u User) (User, error)
	GetAll() ([]User, error)
	GetByID(id int) (User, error)
	Update(id int, req UpdateUserRequest) (User, error)
	Delete(id int) error
}
