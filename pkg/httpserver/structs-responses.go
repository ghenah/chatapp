package httpserver

import (
	"github.com/ghenah/chatapp/pkg/ichatappds"
	"github.com/ghenah/chatapp/pkg/idatastore"
)

type ResponseSuccess struct {
	Success bool `json:"success"`
}

type ResponseAuthSuccess struct {
	idatastore.User
	AccessToken string `json:"accessToken"`
}

type ResponseWSTicket struct {
	WSTicket string `json:"wsTicket"`
}

type ResponseAuthorizedUserInfo struct {
	idatastore.User
}

type ResponseUserSearch struct {
	UsersList []idatastore.UserShort `json:"usersList"`
}

type ResponseChatRoomSearch struct {
	ChatRoomsList []ichatappds.ChatRoomShort `json:"chatRoomsList"`
}
