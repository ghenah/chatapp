package dsgorm

import (
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
