export default {
  saveUserSession(state, userSession) {
    state.loggedIn = true;
    state.accessToken = userSession.accessToken;
    state.id = userSession.id;
    state.username = userSession.username;
    state.email = userSession.email;
    state.friendsList = userSession.friendsList;
    state.ignoreList = userSession.ignoreList;
  },
  updateAccessToken(state, data) {
    state.accessToken = data.accessToken;
  },
  clearUserSession(state) {
    state.loggedIn = false;
    state.accessToken = "";
    state.id = 0;
    state.username = "";
    state.email = "";
    state.friendsList = "";
    state.ignoreList = "";
  },
  setStatusLoggedIn(state) {
    state.loggedIn = true;
  },
  setStatusLoggedOut(state) {
    state.loggedIn = false;
  },
  updateUsername(state, username) {
    state.username = username;
  },
  addFriend(state, friend) {
    if (notInList(state.friendsList, friend)) {
      state.friendsList.push(friend);
    }
  },
  addIgnoredUser(state, ignoredUser) {
    if (notInList(state.ignoreList, ignoredUser)) {
      state.ignoreList.push(ignoredUser);
    }
  },
  removeFriend(state, friendId) {
    state.friendsList = state.friendsList.filter((friend) => {
      return friend.id !== friendId;
    });
  },
  removeIgnoredUser(state, ignoredUserId) {
    state.ignoreList = state.ignoreList.filter((ignoredUser) => {
      return ignoredUser.id !== ignoredUserId;
    });
  },
};

function notInList(array, item) {
  for (let i = 0; i < array.length; i++) {
    if (array[i].id === item.id) {
      return false;
    }
  }

  return true;
}
