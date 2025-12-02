package repositories

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}
