package httpserver

import "github.com/ghenah/chatapp/pkg/idatastore"

type ResponseSuccess struct {
	Success bool `json:"success"`
}

type ResponseAuthSuccess struct {
	idatastore.User
	AccessToken string `json:"access_token"`
}
