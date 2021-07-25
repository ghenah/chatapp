package chatapp

import (
	"fmt"

	"github.com/ghenah/chatapp/pkg/ichatappds"
)

type ChatApp struct {
	InMsgQueue         chan ichatappds.ChatMessage
	clientSessionsList ichatappds.ClientSessionsList
	chatRoomsList      ichatappds.ChatRoomsList
}

func (ca *ChatApp) Start() {
	for {
		select {
		case message := <-ca.InMsgQueue:
			room, err := ca.chatRoomsList.GetRoomInfo(message.RoomID)
			if err != nil {
				fmt.Println("ERROR: ca.chatRoomsList.GetRoomInfo: ", err.Error())

				continue
			}
			if _, present := room.ActiveUsers[message.AuthorID]; !present {
				continue
			}

			// Do NOT broadcast the clientSessionID
			message.ClientSessionID = 0

			// Broadcast the new message to active users in the room

			// Extracting IDs of all room's active users
			activeUsers := make([]uint, len(room.ActiveUsers))
			i := 0
			for k := range room.ActiveUsers {
				activeUsers[i] = k
				i++
			}
			outChannels, err := ca.clientSessionsList.GetOutChannels(activeUsers)
			if err != nil {
				continue
			}
			ca.Broadcast(outChannels, ichatappds.NewChatMessage{
				// 'S' mean security
				Type:           ichatappds.NEW_CHAT_MESSAGE,
				RoomID:         message.RoomID,
				AuthorID:       message.AuthorID,
				AuthorUsername: message.AuthorUsername,
				Content:        message.Content,
			})
		}
	}
}

func (ca *ChatApp) RegisterClientSession(userID uint, username string, outCh chan interface{}) (uint64, error) {
	sessionID, err := ca.clientSessionsList.AddSession(userID, username, outCh)
	if err != nil {
		return 0, err
	}

	return sessionID, nil
}

func (ca *ChatApp) EraseClientSession(sessionID uint64) {
	ca.clientSessionsList.RemoveSessionsByID([]uint64{sessionID})
}

func (ca *ChatApp) GetUserRoomsInfo(userID uint) {
	userRoomsInfo, err := ca.chatRoomsList.GetUserRoomsInfo(userID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	outChannels, err := ca.clientSessionsList.GetOutChannels([]uint{userID})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ca.Broadcast(outChannels, ichatappds.UpdateUserRoomsInfo{
		Type:      ichatappds.UPDATE_USER_ROOMS_INFO,
		RoomsList: userRoomsInfo,
	})
}

// CreateChatRoom creates a chat room specifying the request author as the
// room's owner.
//
// The function does not check whether the user exists within the system. It
// relies on the caller to ensure that and to authenticate the user.
//
// The function returns an error.
func (ca *ChatApp) CreateChatRoom(userID uint, username string, roomVisibility uint, name string) error {
	// Create a room (the room author needn't send a separate request
	// to join the room).
	roomID, activeUsers, err := ca.chatRoomsList.NewRoom(userID, username, roomVisibility, name)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// "Broadcast" the new room to all user's clients.
	outChannels, err := ca.clientSessionsList.GetOutChannels([]uint{userID})
	if err != nil {
		return err
	}
	ca.Broadcast(outChannels, ichatappds.NotificationNewRoom{
		Type:        ichatappds.NOTIFICATION_NEW_ROOM,
		Name:        name,
		ID:          roomID,
		ActiveUsers: activeUsers,
	})

	return nil
}

// GetAllChatRooms returns a list of rooms. The list can be filtered to remove
// rooms belonging to ignored users, if the corresponding flag is set.
func (ca *ChatApp) GetAllChatRooms(userID uint, username string, filterOutIgnored bool) ([]ichatappds.ChatRoomShort, error) {
	roomsInfo, err := ca.chatRoomsList.GetAllRoomsInfoShort()
	if err != nil {
		return nil, err
	}

	if !filterOutIgnored {
		return roomsInfo, nil
	}

	user, err := ds.GetUser(username)
	if err != nil {
		return nil, err
	}
	userActiveRooms, err := ca.chatRoomsList.GetUserRoomsInfo(userID)
	if err != nil {
		return nil, err
	}

	// The filter loop will need to check whether the room owner's ID is in the
	// user's blacklist for each room. Thus, the ignored users' IDs are extracted
	// into a map for faster checks.
	ignoreList := make(map[uint]struct{}, len(user.IgnoreList))
	for _, i := range user.IgnoreList {
		ignoreList[i.ID] = struct{}{}
	}
	// The check is also going to be performed for the room ID, so that the
	// user's active rooms are not returned either.
	activeRooms := make(map[uint]struct{}, len(userActiveRooms))
	for _, i := range userActiveRooms {
		activeRooms[i.ID] = struct{}{}
	}

	// In the best case the resulting list will have the same number of items
	roomsInfoOut := make([]ichatappds.ChatRoomShort, 0, len(roomsInfo))
	// Filter out:
	// - ignored users' rooms
	// - current user's active rooms (including the ones they do not own)
	for _, room := range roomsInfo {
		fmt.Printf("Room ID: %v | Name: %v | Owner: %v\n", room.ID, room.Name, room.OwnerUsername)

		if _, present := ignoreList[room.OwnerID]; present {
			continue
		} else if _, present := activeRooms[room.ID]; present {
			continue
		}

		roomsInfoOut = append(roomsInfoOut, room)
	}

	return roomsInfoOut, nil
}

func (ca *ChatApp) JoinChatRoom(userID uint, username string, roomID uint) {
	err := ca.chatRoomsList.AddUserToRoom(userID, username, roomID)
	if err != nil {
		return
	}

	room, err := ca.chatRoomsList.GetRoomInfo(roomID)
	if err != nil {
		return
	}

	// "Broadcast" the new room to all user's clients (the user who just joined).
	userOutChannel, err := ca.clientSessionsList.GetOutChannels([]uint{userID})
	if err != nil {
		return
	}
	ca.Broadcast(userOutChannel, ichatappds.NotificationNewRoom{
		Type:        ichatappds.NOTIFICATION_NEW_ROOM,
		Name:        room.Name,
		ID:          room.ID,
		ActiveUsers: room.ActiveUsers,
	})

	// The room's active users must be updated with the "new user" notification

	// Extracting IDs of all room's active users
	activeUsers := make([]uint, len(room.ActiveUsers))
	i := 0
	for k := range room.ActiveUsers {
		activeUsers[i] = k
		i++
	}
	outChannels, err := ca.clientSessionsList.GetOutChannels(activeUsers)
	if err != nil {
		return
	}
	ca.Broadcast(outChannels, ichatappds.NotificationUserJoinedRoom{
		Type:     ichatappds.NOTIFICATION_USER_JOINED_ROOM,
		RoomID:   room.ID,
		UserID:   userID,
		Username: username,
	})

}

func (ca *ChatApp) LeaveChatRoom(userID uint, roomID uint) {
	// Grab active users beforehand so the the notification is also sent
	// to the user leaving the room.
	room, err := ca.chatRoomsList.GetRoomInfo(roomID)
	if err != nil {
		return
	}

	// Remove the user from the room's list of active users
	err = ca.chatRoomsList.RemoveUserFromRoom(userID, roomID)
	if err != nil {
		return
	}

	// Broadcast the "user left" notification

	// Extracting IDs of all room's active users
	activeUsers := make([]uint, len(room.ActiveUsers))
	i := 0
	for k := range room.ActiveUsers {
		activeUsers[i] = k
		i++
	}
	outChannels, err := ca.clientSessionsList.GetOutChannels(activeUsers)
	if err != nil {
		return
	}
	ca.Broadcast(outChannels, ichatappds.NotificationUserLeftRoom{
		Type:   ichatappds.NOTIFICATION_USER_LEFT_ROOM,
		RoomID: room.ID,
		UserID: userID,
	})
}

func (ca *ChatApp) Broadcast(outChannels []ichatappds.ClSessChannelWithID, message interface{}) {
	deadChannelsList := []uint64{}

	for _, c := range outChannels {
		select {
		case c.OutCh <- message:
		default:
			close(c.OutCh)
			// Add a dead channel to the list so they could be removed in one
			// transaction after the loop exits
			deadChannelsList = append(deadChannelsList, c.ID)
		}

		if len(deadChannelsList) > 0 {
			ca.clientSessionsList.RemoveSessionsByID(deadChannelsList)
		}
	}
}

func makeIDGenerator() func() uint {
	var currID uint = 1

	return func() uint {
		newID := currID
		currID++
		return newID
	}
}

func makeIDGenerator64() func() uint64 {
	var currID uint64 = 1

	return func() uint64 {
		newID := currID
		currID++
		return newID
	}
}
