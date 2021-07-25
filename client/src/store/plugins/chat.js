const NOTIFICATION_NEW_ROOM = 1;
const NEW_CHAT_MESSAGE = 2;
const UPDATE_USER_ROOMS_INFO = 3;
const NOTIFICATION_USER_JOINED_ROOM = 4;
const NOTIFICATION_USER_LEFT_ROOM = 5;

export default function createWebSocketPlugin() {
  return (store) => {
    store.subscribe((mutation) => {
      if (mutation.type === "chat/initWSConn") {
        mutation.payload.onopen = () => {
          store.commit("chat/setConnected");
          store.dispatch("chat/fetchUsersRooms");
        };
        mutation.payload.onclose = (e) => {
          store.commit("chat/setDisconnected");
          console.log("ws connection closed");
          console.dir(e);
        };
        mutation.payload.onerror = (e) => {
          console.log("ws connection error");
          console.dir(e);
        };
        mutation.payload.onmessage = (e) => {
          try {
            let message = JSON.parse(e.data);
            switch (message.type) {
              case NOTIFICATION_NEW_ROOM:
                store.dispatch("chat/addRoom", {
                  id: message.id,
                  name: message.name,
                  activeUsers: message.activeUsers,
                });

                break;
              case NEW_CHAT_MESSAGE:
                message.type = "newChatMessage";
                store.dispatch("chat/newMessage", message);

                break;
              case UPDATE_USER_ROOMS_INFO:
                store.dispatch("chat/updateUserRoomsInfo", message.roomsList);

                break;

              case NOTIFICATION_USER_JOINED_ROOM:
                store.dispatch("chat/addUserToChat", {
                  type: "userJoined",
                  roomId: message.roomId,
                  userId: message.userId,
                  username: message.username,
                });

                break;
              case NOTIFICATION_USER_LEFT_ROOM:
                message.type = "userLeft";
                store.dispatch("chat/removeUserFromChat", message);

                break;
              default:
                console.log("Error: websocket message: invalid type");
            }
          } catch (err) {
            console.log(
              "Error: websocket message parsing to JSON: " + err.message
            );
            console.dir(e);
            console.log("message data: " + e.data);
          }
        };
      }
    });
  };
}
