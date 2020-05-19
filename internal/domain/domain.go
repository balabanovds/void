package domain

// Storage ...
type Storage interface {
	Users()
}

// UserRepo is a contarct for user reposotories implementations
type UserRepo interface {
	Create(email, password string) error
}
