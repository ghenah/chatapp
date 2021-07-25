package ichatappds

type ChatRoom struct {
	ID            uint            `json:"id"`
	Name          string          `json:"name"`
	OwnerID       uint            `json:"ownerId"`
	OwnerUsername string          `json:"ownerUsername"`
	Visibility    uint            `json:"visibility"`
	ActiveUsers   map[uint]string `json:"activeUsers"`
	PendingUsers  map[uint]struct{}
}

type ChatRoomShort struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	OwnerID       uint   `json:"ownerId"`
	OwnerUsername string `json:"ownerUsername"`
}

const (
	VISIBILITY_PUBLIC  uint = 1
	VISIBILITY_PRIVATE uint = 2
)
