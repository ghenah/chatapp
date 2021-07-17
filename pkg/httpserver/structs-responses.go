package httpserver

import "github.com/ghenah/chatapp/pkg/idatastore"

type ResponseSuccess struct {
	Success bool `json:"success"`
}

type ResponseAuthSuccess struct {
	idatastore.User
	AccessToken string `json:"accessToken"`
}

type ResponseAuthorizedUserInfo struct {
	idatastore.User
}

type ResponseUserSearch struct {
	UsersList []idatastore.UserShort `json:"usersList"`
}
