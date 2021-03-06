export default {
  ticket(state) {
    return state.ticket;
  },
  chatMessages(state) {
    return state.chatMessages;
  },
  wsConn(state) {
    return state.wsConn;
  },
  connected(state) {
    return state.connected;
  },
  messagesLog(state) {
    return state.messageLogs.get(state.activeRoomId);
  },
  roomsList(state) {
    return Array.from(state.roomsList, (room) => room[1]);
  },
  userProfilePics(state) {
    return state.userProfilePics;
  },
  notifications(state) {
    return state.notifications;
  },
  roomVisPublic(state) {
    return state.roomVisPublic;
  },
  roomVisPrivate(state) {
    return state.roomVisPrivate;
  },
  chatRoomSearchResults(state) {
    return state.chatRoomSearchResults;
  },
  activeChatRoomMessages(state) {
    return state.messages.get(state.activeRoomId);
  },
  activeRoomId(state) {
    return state.activeRoomId;
  },
};
