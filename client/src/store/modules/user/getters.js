export default {
  userId(state) {
    return state.id;
  },
  username(state) {
    return state.username;
  },
  email(state) {
    return state.email;
  },
  friendsList(state) {
    return state.friendsList;
  },
  ignoreList(state) {
    return state.ignoreList;
  },
  loggedIn(state) {
    return state.loggedIn;
  },
  accessToken(state) {
    return state.accessToken;
  },
};
