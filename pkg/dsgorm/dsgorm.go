package dsgorm

import (
	"fmt"

	"github.com/ghenah/chatapp/pkg/idatastore"
	"gorm.io/gorm"
)

type DataStoreGORM struct {
	db *gorm.DB
}

// GetDataStore returns the pointer to teh implementation of IDataStore.
func GetDataStore() *DataStoreGORM {
	return dataStore
}

// CreateUser creates a new user entry in the data storage. Unique and valid
// username and email address must be provided. Returns an error.
func (ds *DataStoreGORM) CreateUser(username, email, password string) error {
	userData := User{Username: username, Email: email, Password: password}

	result := ds.db.Create(&userData)
	if result.Error != nil {
		errorMsg := []byte(result.Error.Error())
		switch {
		case mySQLErrors["duplicate entry"].Match(errorMsg):
			return idatastore.ErrorDuplicateEntry
		default:
			return result.Error
		}
	}

	return nil
}

// GetUser returns a user corresponding to the provided username and an error.
func (ds *DataStoreGORM) GetUser(username string) (idatastore.User, error) {
	userResult := User{}
	result := ds.db.Where("username = ?", username).Preload("Friends").Preload("Ignored").Find(&userResult)
	if result.Error != nil {
		fmt.Println(result.Error)
		return idatastore.User{}, result.Error
	}
	if result.RowsAffected == 0 {
		return idatastore.User{}, idatastore.ErrorUserNotFound
	}

	userOut := idatastore.User{
		ID:       userResult.ID,
		Username: userResult.Username,
		Email:    userResult.Email,
		RegDate:  userResult.RegDate,
	}

	for _, f := range userResult.Friends {
		userOut.FriendsList = append(userOut.FriendsList, idatastore.UserShort{
			ID:       f.ID,
			Username: f.Username,
		})
	}

	for _, f := range userResult.Ignored {
		userOut.IgnoreList = append(userOut.IgnoreList, idatastore.UserShort{
			ID:       f.ID,
			Username: f.Username,
		})
	}

	return userOut, nil
}

// GetUserPassword returns the hashed password of a user and an error
func (ds *DataStoreGORM) GetUserPassword(username string) ([]byte, error) {
	userResult := User{}
	result := ds.db.Where("username = ?", username).Find(&userResult)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, idatastore.ErrorUserNotFound
	}

	return []byte(userResult.Password), nil
}

// GetAllUsers returns the list of all users within the system. Only the
// usernames are provided.
func (ds *DataStoreGORM) GetAllUsers() ([]string, error) {
	usersResult := []User{}
	result := ds.db.Order("username asc").Find(&usersResult)
	if result.Error != nil {
		return nil, result.Error
	}

	usersList := []string{}
	for _, u := range usersResult {
		usersList = append(usersList, u.Username)
	}

	return usersList, nil
}

// AddFriend adds a user to the friends list of the current user. Returns
// an error
func (ds *DataStoreGORM) AddFriend(userID, friendID uint) error {
	user := User{
		ID: userID,
	}

	// If the new friend is present in the ignore list, abort the operation
	ignoredUsers := []*User{}
	err := ds.db.Model(&user).Association("Ignored").Find(&ignoredUsers, &User{ID: friendID})
	if err != nil {
		return err
	}
	if len(ignoredUsers) > 0 {
		if ignoredUsers[0].ID == friendID {
			return idatastore.ErrorUserInIgnoreList
		}
	}

	err = ds.db.Model(&user).Association("Friends").Append(&User{ID: friendID})
	if err != nil {
		return err
	}

	return nil
}

// RemoveFriend removes a user from friends list of the current user.
// Returns an error.
func (ds *DataStoreGORM) RemoveFriend(userID, friendID uint) error {
	user := User{
		ID: userID,
	}
	err := ds.db.Model(&user).Association("Friends").Delete(&User{ID: friendID})
	if err != nil {
		return err
	}

	return nil
}

// AddIgnored adds a user to the ignored list of the current user. Returns
// an error
func (ds *DataStoreGORM) AddIgnored(userID, friendID uint) error {
	user := User{
		ID: userID,
	}

	ds.db.Model(&user).Association("Friends").Delete(&User{ID: friendID})
	err := ds.db.Model(&user).Association("Ignored").Append(&User{ID: friendID})
	if err != nil {
		return err
	}

	return nil
}

// RemoveIgnored removes a user from ignored list of the current user.
// Returns an error.
func (ds *DataStoreGORM) RemoveIgnored(userID, friendID uint) error {
	user := User{
		ID: userID,
	}
	err := ds.db.Model(&user).Association("Ignored").Delete(&User{ID: friendID})
	if err != nil {
		return err
	}

	return nil
}
