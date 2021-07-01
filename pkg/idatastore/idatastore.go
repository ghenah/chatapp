package idatastore

type IDataStore interface {
	// CreateUser creates a new user entry in the data storage. Unique and valid
	// username and email address must be provided. Returns an error.
	CreateUser(username, email, password string) error
}
