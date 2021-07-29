package httpserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ghenah/chatapp/pkg/ichatappds"
	"github.com/gorilla/websocket"
)

var (
	appWsOriginSchema string
	appWsOriginDomain string
	appWsOriginPort   string

	upgrader websocket.Upgrader
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type messageDec struct {
	Method string          `json:"method"`
	Data   json.RawMessage `json:"data"`
}
type mdNewMessage struct {
	RoomID  uint   `json:"roomId"`
	Content string `json:"content"`
}
type mdNewChatRoom struct {
	Visibility uint   `json:"visibility"`
	Name       string `json:"name"`
}
type mdJoinChatRoom struct {
	RoomID uint `json:"roomId"`
}
type mdLeaveChatRoom struct {
	RoomID uint `json:"roomId"`
}
type mdInviteUser struct {
	InviteeID       uint   `json:"inviteeId"`
	InviteeUsername string `json:"inviteeUsername"`
	RoomID          uint   `json:"roomId"`
}
type mdStartPersonalChat struct {
	InviteeID       uint   `json:"inviteeId"`
	InviteeUsername string `json:"inviteeUsername"`
}
type mdAcceptInvitation struct {
	RoomId   uint `json:"roomId"`
	Accepted bool `json:"accepted"`
}

type WSConnectionHandler struct {
	ID       uint64
	userID   uint
	username string
	conn     *websocket.Conn
	in       chan ichatappds.ChatMessage
	out      chan interface{}
}

func (h *WSConnectionHandler) read() {
	defer func() {
		ca.EraseClientSession(h.ID)
		h.conn.Close()
	}()

	h.conn.SetReadLimit(maxMessageSize)
	h.conn.SetReadDeadline(time.Now().Add(pongWait))
	h.conn.SetPongHandler(func(string) error {
		h.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		messageDec := &messageDec{}
		err := h.conn.ReadJSON(messageDec)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		fmt.Printf("\nws: user: %s | method: %s\n", h.username, messageDec.Method) // RRemove
		fmt.Printf("ws: dataNonDecoded: %s\n\n", messageDec.Data)                  // RRemove

		switch messageDec.Method {
		case "newMessage":
			data := &mdNewMessage{}
			err = json.Unmarshal(messageDec.Data, data)
			if err != nil {
				fmt.Println("unmarshalling newMessage.Data: ", err.Error())
				h.conn.Close()
			}

			h.in <- ichatappds.ChatMessage{
				ClientSessionID: h.ID,
				RoomID:          data.RoomID,
				Content:         data.Content,
				AuthorID:        h.userID,
				AuthorUsername:  h.username,
			}
		case "newChatRoom":
			data := &mdNewChatRoom{}
			err = json.Unmarshal(messageDec.Data, data)
			if err != nil {
				fmt.Println("unmarshalling newChatRoom.Data: ", err.Error())
				h.conn.Close()

				continue
			}

			go ca.CreateChatRoom(h.userID, h.username, data.Visibility, data.Name)

		case "userRoomsInfo":
			go ca.GetUserRoomsInfo(h.userID)
		case "joinChatRoom":
			data := &mdJoinChatRoom{}
			err = json.Unmarshal(messageDec.Data, data)
			if err != nil {
				fmt.Println("unmarshalling joinChatRoom.Data: ", err.Error())
				h.conn.Close()

				continue
			}

			go ca.JoinChatRoom(h.userID, h.username, data.RoomID)
		case "leaveChatRoom":
			data := &mdLeaveChatRoom{}
			err = json.Unmarshal(messageDec.Data, data)
			if err != nil {
				fmt.Println("unmarshalling leaveChatRoom.Data: ", err.Error())
				h.conn.Close()

				continue
			}

			go ca.LeaveChatRoom(h.userID, data.RoomID)
		case "inviteUser":
			data := &mdInviteUser{}
			err = json.Unmarshal(messageDec.Data, data)
			if err != nil {
				fmt.Println("unmarshalling inviteUser.Data: ", err.Error())
				h.conn.Close()

				continue
			}

			go ca.InviteUser(h.userID, data.InviteeID, data.RoomID, h.username, data.InviteeUsername)
		case "acceptInvitation":
			data := &mdAcceptInvitation{}
			err = json.Unmarshal(messageDec.Data, data)
			if err != nil {
				fmt.Println("unmarshalling acceptInvitation.Data: ", err.Error())
				h.conn.Close()

				continue
			}

			go ca.AcceptInvitation(h.userID, h.username, data.RoomId, data.Accepted)
		case "startPersonalChat":
			data := &mdStartPersonalChat{}
			err = json.Unmarshal(messageDec.Data, data)
			if err != nil {
				fmt.Println("unmarshalling startPersonalChat.Data: ", err.Error())
				h.conn.Close()

				continue
			}

			go ca.StartPersonalChat(h.userID, h.username, data.InviteeID, data.InviteeUsername)
		default:
			fmt.Println("ws connection handler: non-protocol message")
		}

	}
}

func (h *WSConnectionHandler) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		h.conn.Close()
	}()

	for {
		select {
		case message, ok := <-h.out:
			h.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The channel has been closed
				h.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := h.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			wsOutMessage := []byte{}

			// Send the current message.
			fmt.Printf("User: %s, ClConnID: %d, Message: %v\n", h.username, h.ID, message) // RRemove

			switch message.(type) {
			case ichatappds.NewChatMessage, ichatappds.NotificationNewRoom, ichatappds.UpdateUserRoomsInfo, ichatappds.NotificationUserJoinedRoom, ichatappds.NotificationUserLeftRoom, ichatappds.NotificationRoomInvitation, ichatappds.NotificationRoomDeleted:
				wsOutMessage, err = json.Marshal(message)
				if err != nil {
					continue
				}
			}
			w.Write(wsOutMessage)
			if err := w.Close(); err != nil {
				return
			}

			// Send the messages waiting in the queue.
			n := len(h.out)
			for i := 0; i < n; i++ {
				w, err := h.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}

				message = <-h.out
				switch message.(type) {
				case ichatappds.ChatMessage:
					wsOutMessage, err = json.Marshal(message)
					if err != nil {
						continue
					}
				}
				w.Write(wsOutMessage)
				if err := w.Close(); err != nil {
					return
				}
			}
		case <-ticker.C:
			h.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := h.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func serveConnection(w http.ResponseWriter, r *http.Request, in chan ichatappds.ChatMessage, userID uint, username string) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	outCh := make(chan interface{}, 64)

	sessionID, err := ca.RegisterClientSession(userID, username, outCh)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	wsConnHandler := &WSConnectionHandler{
		ID:       sessionID,
		userID:   userID,
		username: username,
		conn:     conn,
		in:       in,
		out:      outCh,
	}

	go wsConnHandler.read()
	go wsConnHandler.write()

	return nil
}
