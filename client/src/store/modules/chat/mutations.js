export default {
  saveTicket(state, ticket) {
    state.ticket = ticket;
  },
  clearTicket(state) {
    state.ticket = "";
  },
  initWSConn(state, wsConn) {
    state.wsConn = wsConn;
  },
  setConnected(state) {
    state.connected = true;
  },
  setDisconnected(state) {
    state.connected = false;
    state.roomsList = new Map();
    state.messageLogs = new Map();
    state.conn = null;
    state.ticket = "";
    state.chatRoomSearchResults = [];
  },
  addRoom(state, room) {
    state.roomsList.set(room.id, room);

    // Initialize the chat room messages list in a separate structure
    state.messageLogs.set(room.id, Array());
  },
  deleteRoom(state, roomId) {
    state.roomsList.delete(roomId);
    state.messageLogs.delete(roomId);
  },
  removeChatRoom(state, roomId) {
    state.roomsList.delete(roomId);
    state.messageLogs.delete(roomId);
  },
  addNewMessage(state, d) {
    let roomMsgLog = state.messageLogs.get(d.roomId);
    roomMsgLog.push(d);
    state.messageLogs.set(d.roomId, roomMsgLog);
  },
  addUserToChat(state, d) {
    let room = state.roomsList.get(d.roomId);
    room.activeUsers[d.userId] = d.username;
    state.roomsList.set(d.roomId, room);
  },
  addNotificationInvitation(state, d) {
    d.id = generateNotificationId();
    state.notifications.push(d);
  },
  removeUserFromChat(state, d) {
    let room = state.roomsList.get(d.roomId);
    room.activeUsers[d.userId];
    delete room.activeUsers[d.userId];
    state.roomsList.set(d.roomId, room);
  },
  addChatLog(state, d) {
    // Get the username to use in the message log
    let room = state.roomsList.get(d.roomId);
    let username = room.activeUsers[d.userId];
    d.username = username;

    let roomMsgLog = state.messageLogs.get(d.roomId);
    roomMsgLog.push(d);
    state.messageLogs.set(d.roomId, roomMsgLog);
  },
  setActiveChatRoom(state, roomId) {
    state.activeRoomId = roomId;
  },
  updateUserRoomsInfo(state, roomsList) {
    state.roomsList = roomsList;

    // Make a map of empty arrays, based on elements of the source map
    let easyMapper = Array.from(roomsList, (value) => value[1]);
    state.messageLogs = new Map(easyMapper.map((e) => [e.id, Array({})]));
  },
  removeNotification(state, invId) {
    console.log("state.notifications:");
    console.dir(state.notifications);
    console.log("invId: " + invId);
    state.notifications = state.notifications.filter((e) => e.id !== invId);
  },
  saveChatRoomSearchResults(state, chatRoomSearchResults) {
    state.chatRoomSearchResults = chatRoomSearchResults;
  },
  cleanUpSessionInfo(state) {
    state.roomsList = new Map();
    state.messageLogs = new Map();
    state.notifications = [];
    state.connected = false;
    state.ticket = "";
    state.chatRoomSearchResults = [];
    if (state.wsConn !== null) {
      state.wsConn.close();
    }
    state.wsConn = null;
    console.log("WSCONN DOWN");
  },
};

const generateNotificationId = (function() {
  let id = 1;
  return function() {
    return id++;
  };
})();
