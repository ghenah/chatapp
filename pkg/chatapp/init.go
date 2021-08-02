package chatapp

import (
	"github.com/ghenah/chatapp/pkg/ichatappds"
	"github.com/ghenah/chatapp/pkg/idatastore"
)

type ChatAppConfig struct {
	ClientSessionsList ichatappds.ClientSessionsList
	ChatRoomsList      ichatappds.ChatRoomsList
	ProfilePictures    ichatappds.ProfilePictures
	UsersDS            idatastore.IDataStore
}

var ds idatastore.IDataStore

func Init(cfg *ChatAppConfig) (*ChatApp, error) {
	ds = cfg.UsersDS

	// If no external list management implementations were used
	// provide the built-in ones.
	if cfg.ClientSessionsList == nil {
		cfg.ClientSessionsList = &ClientSessionsList{
			clientSessions:     make(map[uint64]*ichatappds.ClientSession),
			ownerSessionsTable: make(map[uint][]uint64),
		}
	}
	if cfg.ChatRoomsList == nil {
		cfg.ChatRoomsList = &ChatRoomsList{
			chatRooms:           make(map[uint]*ichatappds.ChatRoom),
			userActiveRoomsList: make(map[uint][]uint),
		}
	}
	if cfg.ProfilePictures == nil {
		cfg.ProfilePictures = &ProfilePictures{
			pictures: make(map[uint]string),
		}
	}

	chatApp := &ChatApp{
		InMsgQueue:         make(chan ichatappds.ChatMessage),
		clientSessionsList: cfg.ClientSessionsList,
		chatRoomsList:      cfg.ChatRoomsList,
		profilePictures:    cfg.ProfilePictures,
	}
	go chatApp.Start()

	return chatApp, nil
}
