package chatapp

import (
	"sync"

	"github.com/ghenah/chatapp/pkg/ichatappds"
)

var TEMPORARYgenChatRoomID = makeIDGenerator()

type ChatRoomsList struct {
	sync.Mutex
	// [roomID]*room
	chatRooms map[uint]*ichatappds.ChatRoom
	// [userID][]roomID
	userActiveRoomsList map[uint][]uint
}

func (crl *ChatRoomsList) NewRoom(userID uint, username string, roomVisibility uint, name string) (uint, map[uint]string, error) {
	switch roomVisibility {
	case ichatappds.VISIBILITY_PRIVATE, ichatappds.VISIBILITY_PUBLIC, ichatappds.VISIBILITY_PERSONAL:
		room := &ichatappds.ChatRoom{
			ID:            TEMPORARYgenChatRoomID(),
			Name:          name,
			OwnerID:       userID,
			OwnerUsername: username,
			Visibility:    roomVisibility,
			ActiveUsers:   map[uint]string{userID: username},
			PendingUsers:  map[uint]struct{}{},
		}

		crl.Lock()
		crl.chatRooms[room.ID] = room
		crl.userActiveRoomsList[userID] = append(crl.userActiveRoomsList[userID], room.ID)
		crl.Unlock()

		return room.ID, room.ActiveUsers, nil
	default:
	}

	return 0, nil, ichatappds.ErrorInvalidChatVisibility
}

func (crl *ChatRoomsList) AddUserToRoom(userID uint, username string, roomID uint) error {
	crl.Lock()
	defer crl.Unlock()
	room, ok := crl.chatRooms[roomID]
	if !ok {
		return ichatappds.ErrorChatRoomDoesNotExist
	}

	// If the room is private, check whether the user has been invited
	if room.Visibility == ichatappds.VISIBILITY_PRIVATE || room.Visibility == ichatappds.VISIBILITY_PERSONAL {
		if _, present := room.PendingUsers[userID]; !present {
			return ichatappds.ErrorUserIsNotInvited
		}
	}

	room.ActiveUsers[userID] = username
	// If a user has joined the room by accepting an invitation, clear their
	// invitee status.
	delete(room.PendingUsers, userID)

	crl.userActiveRoomsList[userID] = append(crl.userActiveRoomsList[userID], roomID)

	return nil
}

func (crl *ChatRoomsList) RemoveUserFromRoom(userID, roomID uint) error {
	crl.Lock()
	defer crl.Unlock()
	room, ok := crl.chatRooms[roomID]
	if !ok {
		return ichatappds.ErrorChatRoomDoesNotExist
	}

	delete(room.ActiveUsers, userID)
	userActiveRooms := crl.userActiveRoomsList[userID]
	for i, e := range userActiveRooms {
		if e == roomID {
			crl.userActiveRoomsList[userID] = append(userActiveRooms[:i], userActiveRooms[i+1:]...)

			break
		}
	}

	return nil
}

func (crl *ChatRoomsList) RemoveUserFromPending(userID, roomID uint) error {
	crl.Lock()
	defer crl.Unlock()
	room, ok := crl.chatRooms[roomID]
	if !ok {
		return ichatappds.ErrorChatRoomDoesNotExist
	}

	delete(room.PendingUsers, userID)
	return nil
}

func (crl *ChatRoomsList) GetRoomInfo(roomID uint) (ichatappds.ChatRoom, error) {
	roomOut := ichatappds.ChatRoom{}
	activeUsers := map[uint]string{}

	crl.Lock()
	defer crl.Unlock()
	room, ok := crl.chatRooms[roomID]
	if !ok {
		return ichatappds.ChatRoom{}, ichatappds.ErrorChatRoomDoesNotExist
	}

	// deep copying active users
	for uID, un := range room.ActiveUsers {
		activeUsers[uID] = un
	}

	roomOut = ichatappds.ChatRoom{
		ID:            room.ID,
		Name:          room.Name,
		OwnerID:       room.OwnerID,
		OwnerUsername: room.OwnerUsername,
		Visibility:    room.Visibility,
		ActiveUsers:   activeUsers,
	}

	return roomOut, nil
}

func (crl *ChatRoomsList) GetRoomName(roomID uint) (string, error) {
	crl.Lock()
	defer crl.Unlock()

	if room, ok := crl.chatRooms[roomID]; ok {
		return room.Name, nil
	}

	return "", ichatappds.ErrorChatRoomDoesNotExist
}

func (crl *ChatRoomsList) GetUserRoomsInfo(userID uint) ([]ichatappds.ChatRoom, error) {
	crl.Lock()
	defer crl.Unlock()
	userRoomsInfo := []ichatappds.ChatRoom{}
	for _, roomID := range crl.userActiveRoomsList[userID] {
		room := crl.chatRooms[roomID]
		activeUsers := map[uint]string{}
		// deep copying active users
		for uID, un := range room.ActiveUsers {
			activeUsers[uID] = un
		}
		userRoomsInfo = append(userRoomsInfo, ichatappds.ChatRoom{
			ID:            room.ID,
			Name:          room.Name,
			OwnerID:       room.OwnerID,
			OwnerUsername: room.OwnerUsername,
			Visibility:    room.Visibility,
			ActiveUsers:   activeUsers,
		})
	}

	return userRoomsInfo, nil
}

// GetAllRoomsInfoShort returns a list of active chat rooms. The private rooms are
// not included in the list.
func (crl *ChatRoomsList) GetAllRoomsInfoShort() ([]ichatappds.ChatRoomShort, error) {
	roomsOut := []ichatappds.ChatRoomShort{}
	crl.Lock()
	defer crl.Unlock()
	for _, room := range crl.chatRooms {
		if room.Visibility == ichatappds.VISIBILITY_PRIVATE {
			continue
		}

		roomsOut = append(roomsOut, ichatappds.ChatRoomShort{
			ID:            room.ID,
			Name:          room.Name,
			OwnerID:       room.OwnerID,
			OwnerUsername: room.OwnerUsername,
		})
	}

	return roomsOut, nil
}

func (crl *ChatRoomsList) AddInvitee(userID, inviteeID, roomID uint) error {
	crl.Lock()
	defer crl.Unlock()
	room, ok := crl.chatRooms[roomID]
	if !ok {
		return ichatappds.ErrorChatRoomDoesNotExist
	}
	// The current logic is that the invitation cannot be resent repeatedly or
	// sent when the user is already an active member of the chat room. Thus,
	// this check is performed as the first one, regardless of the
	// room/ownership statuses.
	if _, present := room.PendingUsers[inviteeID]; present {
		return ichatappds.ErrorUserAlreadyPending
	} else if _, present := room.ActiveUsers[inviteeID]; present {
		return ichatappds.ErrorUserAlreadyActive
	}

	switch room.Visibility {
	case ichatappds.VISIBILITY_PUBLIC:
		room.PendingUsers[inviteeID] = struct{}{}
		return nil
	case ichatappds.VISIBILITY_PRIVATE:
		if userID != room.OwnerID {
			return ichatappds.ErrorUserNotRoomOwner
		}
	case ichatappds.VISIBILITY_PERSONAL:
		// The app automatically invites the target user to the chat and a personal
		// chat can only have two parties in it.
		if userID != room.OwnerID {
			return ichatappds.ErrorUserNotRoomOwner
		} else if len(room.PendingUsers) > 0 || len(room.ActiveUsers) == 2 {
			return ichatappds.ErrorCannotInviteToPersonal
		}

	}

	room.PendingUsers[inviteeID] = struct{}{}
	return nil
}

func (crl *ChatRoomsList) GetRoomVisibility(roomID uint) (uint, error) {
	crl.Lock()
	defer crl.Unlock()

	room, ok := crl.chatRooms[roomID]
	if !ok {
		return 0, ichatappds.ErrorChatRoomDoesNotExist
	}

	return room.Visibility, nil
}

// DeleteRoom deletes a room from the list and returns a list of
// users that were active at the moment of deletion.
func (crl *ChatRoomsList) DeleteRoom(roomID uint) []uint {
	crl.Lock()
	defer crl.Unlock()

	room, ok := crl.chatRooms[roomID]
	if !ok {
		return []uint{}
	}
	activeUsersIDs := make([]uint, len(room.ActiveUsers))

	for userID := range room.ActiveUsers {
		activeUsersIDs = append(activeUsersIDs, userID)
		userActiveRooms := crl.userActiveRoomsList[userID]
		for i, e := range userActiveRooms {
			if e == roomID {
				crl.userActiveRoomsList[userID] = append(userActiveRooms[:i], userActiveRooms[i+1:]...)

				break
			}
		}
	}

	return activeUsersIDs
}

// removeRoomReference removes the room ID from the list of rooms
// of a particular user (the user must be separately removed from the
// list of active users of the room; the function only removes the
// reference)
func removeRoomReference(s []uint, e uint) []uint {
	sOut := []uint{}
	for _, el := range s {
		if e != el {
			sOut = append(sOut, el)
		}
	}

	return sOut
}
