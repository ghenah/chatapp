import mutations from "./mutations.js";
import getters from "./getters.js";
import actions from "./actions.js";

const CHAT_ROOM_VISIBILITY_PUBLIC = 1;
const CHAT_ROOM_VISIBILITY_PRIVATE = 2;

export default {
  namespaced: true,
  state() {
    return {
      ticket: "",
      wsConn: null,
      connected: false,
      activeRoomId: 0,
      messageLogs: new Map(),
      roomsList: new Map(),
      roomVisPublic: CHAT_ROOM_VISIBILITY_PUBLIC,
      roomVisPrivate: CHAT_ROOM_VISIBILITY_PRIVATE,
      chatRoomSearchResults: [],
    };
  },
  mutations,
  getters,
  actions,
};
