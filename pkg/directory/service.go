package directory

// Service is an interface that implements getting a user
type Service interface {
	GetUser(email string) (*User, error)
}
