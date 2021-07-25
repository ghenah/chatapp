package ichatappds

type ChatRoomsList interface {
	NewRoom(userID uint, username string, chatVisibility uint, name string) (uint, map[uint]string, error)
	GetRoomInfo(roomID uint) (ChatRoom, error)
	GetUserRoomsInfo(userID uint) ([]ChatRoom, error)
	GetAllRoomsInfoShort() ([]ChatRoomShort, error)
	AddUserToRoom(userID uint, username string, roomID uint) error
	RemoveUserFromRoom(userID, roomID uint) error
}
