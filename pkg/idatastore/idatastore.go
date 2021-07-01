package idatastore

type IDataStore interface {
	// CreateUser creates a new user entry in the data storage. Unique and valid
	// username and email address must be provided. Returns an error.
	CreateUser(username, email, password string) error

	// GetUser returns a user corresponding to the provided ID and an error.
	GetUser(username string) (User, error)

	// GetUserPassword returns a hashed password of a user and an error.
	GetUserPassword(username string) ([]byte, error)
}
