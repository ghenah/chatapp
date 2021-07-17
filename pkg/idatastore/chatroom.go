package idatastore

type ChatRoom struct {
	ID           uint        `json:"id"`
	Owner        UserShort   `json:"owner"`
	ActiveUsers  []UserShort `json:"activeUsers"`
	PendingUsers []UserShort `json:"pendingUsers"`
}
