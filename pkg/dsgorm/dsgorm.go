package dsgorm

import (
	"gorm.io/gorm"
)

type DataStoreGORM struct {
	db *gorm.DB
}

// GetDataStore returns the pointer to teh implementation of IDataStore.
func GetDataStore() *DataStoreGORM {
	return dataStore
}
