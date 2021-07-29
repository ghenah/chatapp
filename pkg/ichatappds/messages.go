package ichatappds

type ChatMessage struct {
	AuthorID        uint
	AuthorUsername  string
	ClientSessionID uint64
	RoomID          uint
	Content         string
}

type NewChatMessage struct {
	Type           uint8  `json:"type"`
	AuthorID       uint   `json:"authorId"`
	AuthorUsername string `json:"authorUsername"`
	RoomID         uint   `json:"roomId"`
	Content        string `json:"content"`
}

type NotificationNewRoom struct {
	Type        uint8           `json:"type"`
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Visibility  uint            `json:"visibility"`
	ActiveUsers map[uint]string `json:"activeUsers"`
}

type NotificationUserJoinedRoom struct {
	Type     uint8  `json:"type"`
	RoomID   uint   `json:"roomId"`
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
}

type NotificationUserLeftRoom struct {
	Type   uint8 `json:"type"`
	RoomID uint  `json:"roomId"`
	UserID uint  `json:"userId"`
}

type NotificationRoomDeleted struct {
	Type   uint8 `json:"type"`
	RoomID uint  `json:"roomId"`
}

type NotificationRoomInvitation struct {
	Type       uint8  `json:"type"`
	RoomID     uint   `json:"roomId"`
	UserID     uint   `json:"userId"`
	Username   string `json:"username"`
	RoomName   string `json:"roomName"`
	Visibility uint   `json:"visibility"`
}

type UpdateUserRoomsInfo struct {
	Type      uint8      `json:"type"`
	RoomsList []ChatRoom `json:"roomsList"`
}

const (
	NOTIFICATION_NEW_ROOM         uint8 = 1
	NEW_CHAT_MESSAGE              uint8 = 2
	UPDATE_USER_ROOMS_INFO        uint8 = 3
	NOTIFICATION_USER_JOINED_ROOM uint8 = 4
	NOTIFICATION_USER_LEFT_ROOM   uint8 = 5
	NOTIFICATION_ROOM_INVITATION  uint8 = 6
	NOTIFICATION_ROOM_DELETED     uint8 = 7
)
