package ichatappds

import "errors"

var (
	ErrorInvalidChatVisibility = errors.New("invalid chat visibility")
	ErrorChatRoomDoesNotExist  = errors.New("chat room does not exist")
	ErrorUserIsNotInvited      = errors.New("user is not invited")
)
