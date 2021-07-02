package idatastore

type IDataStore interface {
	// CreateUser creates a new user entry in the data storage. Unique and valid
	// username and email address must be provided. Returns an error.
	CreateUser(username, email, password string) error

	// GetUser returns a user corresponding to the provided ID and an error.
	GetUser(username string) (User, error)

	// GetUserPassword returns a hashed password of a user and an error.
	GetUserPassword(username string) ([]byte, error)

	// GetAllUsers returns the list of all users within the system. Only the
	// usernames are provided.
	GetAllUsers() ([]string, error)

	// AddFriend adds a user to the friends list of the current user. If the
	// user is in the ignored list, abort the operation. Returns
	// an error.
	AddFriend(userID, friendID uint) error

	// RemoveFriend removes a user from friends list of the current user.
	// Returns an error.
	RemoveFriend(userID, friendID uint) error

	// AddIgnored adds a user to the ignore list of the current user. Returns
	// an error. If the ignored user is in the friends list, they are removed
	// from it.
	AddIgnored(userID, ignoredID uint) error

	// RemoveIgnored removes a user from friends list of the current user.
	// Returns an error. Any former "friends" are not re-added back into the
	// friends list.
	RemoveIgnored(userID, ignoredID uint) error

	// UpdateUserPassword updates the password of the user. Returns an error.
	UpdateUserPassword(userID uint, password string) error

	// UpdateUsername updates the username. Returns an error.
	UpdateUsername(userID uint, newUsername string) error
}
