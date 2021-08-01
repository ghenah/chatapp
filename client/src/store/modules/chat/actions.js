var schema = process.env.VUE_APP_SCHEMA;
var domain = process.env.VUE_APP_DOMAIN;
var port = process.env.VUE_APP_PORT;
var address = schema + domain + ":" + port;

export default {
  getTicket(context) {
    return new Promise((resolve, reject) => {
      context
        .dispatch("user/getAccessToken", null, {
          root: true,
        })
        .then((accessToken) => {
          sendRequest(address + "/api/v1/chat/ticket", "GET", {
            Authorization: "Bearer " + accessToken,
          }).then((response) => {
            if (response.status.ok) {
              context.commit("saveTicket", response.data.wsTicket);

              resolve();
            } else {
              reject(response.data.message);
            }
          });
        })
        .catch((e) => {
          console.log(e);
        });
    });
  },
  openWS(context) {
    let ticket = context.getters["ticket"];
    let url =
      "ws:" +
      "//" +
      process.env.VUE_APP_DOMAIN +
      ":" +
      process.env.VUE_APP_PORT +
      "/ws/connect" +
      "?ticket=" +
      ticket;

    context.commit("initWSConn", new WebSocket(url));
    context.commit("clearTicket");
  },
  fetchUsersRooms(context) {
    let wsConn = context.getters["wsConn"];
    let payload = {
      method: "userRoomsInfo",
    };
    wsConn.send(JSON.stringify(payload));
  },
  searchChatRooms(context) {
    context
      .dispatch("user/getAccessToken", null, {
        root: true,
      })
      .then((accessToken) => {
        sendRequest(address + "/api/v1/chat/rooms/search", "GET", {
          "Content-Type": "application/json",
          Authorization: "Bearer " + accessToken,
        }).then((response) => {
          if (response.status.ok) {
            context.commit(
              "saveChatRoomSearchResults",
              response.data.chatRoomsList
            );
          } else {
            console.log(response.data.message);
          }
        });
      })
      .catch((e) => {
        console.log(e);
      });
  },
  newRoom(context, d) {
    let wsConn = context.getters["wsConn"];
    let payload = {
      method: "newChatRoom",
      data: { visibility: d.visibility, name: d.name },
    };
    wsConn.send(JSON.stringify(payload));
  },
  addRoom(context, room) {
    context.commit("addRoom", room);
  },
  deleteRoom(context, d) {
    context.commit("deleteRoom", d.roomId);
  },
  updateUserProfilePics(context, profilePics) {
    context.commit("updateUserProfilePics", profilePics);
    context.commit("updateUserProfilePicture", {
      userId: context.rootGetters["user/userId"],
      picture: context.rootGetters["user/profilePicture"],
    });
  },
  updateUserRoomsInfo(context, roomsList) {
    context.commit(
      "updateUserRoomsInfo",
      new Map(roomsList.map((i) => [i.id, i]))
    );
  },
  sendMessage(context, d) {
    let wsConn = context.getters["wsConn"];
    let payload = {
      method: "newMessage",
      data: d,
    };
    wsConn.send(JSON.stringify(payload));
  },
  newMessage(context, d) {
    let ignoreList = context.rootGetters["user/ignoreList"];
    for (const e of ignoreList) {
      if (e.id === d.authorId) {
        return;
      }
    }

    context.commit("addNewMessage", d);
  },
  setActiveChatRoom(context, data) {
    context.commit("setActiveChatRoom", data.roomId);
  },
  joinChatRoom(context, data) {
    let wsConn = context.getters["wsConn"];
    let payload = {
      method: "joinChatRoom",
      data: { roomId: data.roomId },
    };
    wsConn.send(JSON.stringify(payload));
  },
  leaveChatRoom(context, d) {
    let wsConn = context.getters["wsConn"];
    let payload = {
      method: "leaveChatRoom",
      data: { roomId: d.roomId },
    };
    wsConn.send(JSON.stringify(payload));
  },
  addUserToChat(context, d) {
    context.commit("updateUserProfilePicture", {
      userId: d.userId,
      picture: d.picture,
    });
    context.commit("addUserToChat", d);
    let ignoreList = context.rootGetters["user/ignoreList"];
    for (const e of ignoreList) {
      if (e.id === d.userId) {
        return;
      }
    }
    context.commit("addChatLog", d);
  },
  removeUserFromChat(context, d) {
    if (d.userId === context.rootGetters["user/userId"]) {
      context.commit("removeChatRoom", d.roomId);

      return;
    }

    let ignoreList = context.rootGetters["user/ignoreList"];
    for (const e of ignoreList) {
      console.dir(e);
      if (e.id === d.userId) {
        return;
      }
    }
    context.commit("addChatLog", d);

    context.commit("removeUserFromChat", d);
  },
  inviteUser(context, d) {
    let wsConn = context.getters["wsConn"];
    let payload = {
      method: "inviteUser",
      data: {
        inviteeId: d.inviteeId,
        roomId: d.roomId,
        inviteeUsername: d.inviteeUsername,
      },
    };
    wsConn.send(JSON.stringify(payload));
  },
  startPersonalChat(context, d) {
    let wsConn = context.getters["wsConn"];
    let payload = {
      method: "startPersonalChat",
      data: {
        inviteeId: d.userId,
        inviteeUsername: d.username,
      },
    };
    wsConn.send(JSON.stringify(payload));
  },
  addNotificationInvitation(context, d) {
    context.commit("addNotificationInvitation", d);
  },
  acceptInvitation(context, d) {
    let wsConn = context.getters["wsConn"];
    let payload = {
      method: "acceptInvitation",
      data: {
        roomId: d.roomId,
        accepted: d.accepted,
      },
    };
    wsConn.send(JSON.stringify(payload));
    context.commit("removeNotification", d.invId);
  },
};

async function sendRequest(url, method, headers, body) {
  let response = await fetch(url, {
    method,
    headers,
    body: JSON.stringify(body),
  });

  let output = {
    status: {
      ok: response.ok,
      code: response.status,
      text: response.statusText,
    },
  };
  output.data = await response.json();

  return output;
}
