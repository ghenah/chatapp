package idatastore

import "errors"

var (
	ErrorUserNotFound     = errors.New("user not found")
	ErrorDuplicateEntry   = errors.New("duplicate entry")
	ErrorUserInIgnoreList = errors.New("user in ignore list")
)
