package repositories

type IUserRepository interface {
}

type UserRepository struct{}

func NewUserRepository() IUserRepository {
	return &UserRepository{}
}
