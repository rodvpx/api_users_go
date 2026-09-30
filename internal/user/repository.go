package user

type Repository interface {
	Create(u User) (User, error)
	GetAll() ([]User, error)
	GetByID(id int) (User, error)
	Update(id int, u User) (User, error)
	Delete(id int) error
}
