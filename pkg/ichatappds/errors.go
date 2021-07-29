package ichatappds

import "errors"

var (
	ErrorInvalidChatVisibility  = errors.New("invalid chat visibility")
	ErrorChatRoomDoesNotExist   = errors.New("chat room does not exist")
	ErrorUserIsNotInvited       = errors.New("user is not invited")
	ErrorUserAlreadyActive      = errors.New("user already active")
	ErrorUserAlreadyPending     = errors.New("user already pending")
	ErrorCannotInviteToPersonal = errors.New("cannot invite to personal")
	ErrorUserNotRoomOwner       = errors.New("user not room owner")
)
