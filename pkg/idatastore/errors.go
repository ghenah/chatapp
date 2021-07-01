package idatastore

import "errors"

var (
	ErrorUserNotFound   = errors.New("user not found")
	ErrorDuplicateEntry = errors.New("duplicate entry")
)
