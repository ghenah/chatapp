import mutations from "./mutations.js";
import getters from "./getters.js";
import actions from "./actions.js";

export default {
  namespaced: true,
  state() {
    return {
      loggedIn: window.localStorage.getItem("loggedIn") || false,
      accessToken: window.localStorage.getItem("accessToken") || "",
      id: 0,
      username: "",
      profilePicture: "",
      email: "",
      friendsList: [],
      ignoreList: [],
    };
  },
  mutations,
  getters,
  actions,
};
