package idatastore

import "time"

type User struct {
	ID          uint        `json:"id"`
	Username    string      `json:"username"`
	Email       string      `json:"email"`
	FriendsList []UserShort `json:"friendsList"`
	IgnoreList  []UserShort `json:"ignoreList"`
	RegDate     time.Time   `json:"regDate"`
}

type UserShort struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}
