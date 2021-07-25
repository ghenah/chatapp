var schema = process.env.VUE_APP_SCHEMA;
var domain = process.env.VUE_APP_DOMAIN;
var port = process.env.VUE_APP_PORT;
var address = schema + domain + ":" + port;

export default {
  login(context, data) {
    return new Promise((resolve, reject) => {
      sendRequest(
        address + "/auth/signin",
        "POST",
        { "Content-Type": "application/json" },
        { username: data.username, password: data.password }
      ).then((response) => {
        if (response.status.ok) {
          context.commit("saveUserSession", response.data);
          // Also persist the login session
          saveToLocalStorage({
            accessToken: response.data.accessToken,
            loggedIn: true,
          });

          // Auto-authenticate the user within the chat (websocket) itself
          context
            .dispatch("chat/getTicket", null, { root: true })
            .then(() => {
              context.dispatch("chat/openWS", null, { root: true });

              resolve();
              return;
            })
            .catch((errorMsg) => {
              reject(errorMsg);
              return;
            });

          resolve();
        } else {
          reject(response.data.message);
        }
      });
    });
  },
  logout(context) {
    context.commit("clearUserSession");
    // Clear the chat app info ()
    context.commit("chat/cleanUpSessionInfo", null, { root: true });
    // Also clear the session-related items from the local storage
    ["loggedIn", "accessToken"].forEach((key) => {
      window.localStorage.removeItem(key);
    });
  },
  setStatusLoggedIn(context) {
    context.commit("setStatusLoggedIn");
  },
  setStatusLoggedOut(context) {
    context.commit("setStatusLoggedOut");
  },
  getProfileInfo(context) {
    return new Promise((resolve, reject) => {
      let accessToken = context.getters["accessToken"];
      sendRequest(address + "/api/v1/users/profile", "GET", {
        "Content-Type": "application/json",
        Authorization: "Bearer " + accessToken,
      }).then((response) => {
        if (response.status.ok) {
          // Access token needs to be added to the response so that the
          // saveUserSession mutation can be used instead of defining a
          // new function
          response.data["accessToken"] = accessToken;
          context.commit("saveUserSession", response.data);
          resolve();
        } else {
          reject(response.data.message);
        }
      });
    });
  },
  addFriend(context, friend) {
    let userId = context.getters["userId"];
    let accessToken = context.getters["accessToken"];
    sendRequest(
      address + "/api/v1/users/friends",
      "POST",
      {
        "Content-Type": "application/json",
        Authorization: "Bearer " + accessToken,
      },
      {
        userId: userId,
        friendId: friend.id,
      }
    ).then((response) => {
      if (response.status.ok) {
        context.commit("addFriend", friend);
      } else {
        console.log(response.data.message);
      }
    });
  },
  ignoreUser(context, ignoredUser) {
    let userId = context.getters["userId"];
    let accessToken = context.getters["accessToken"];
    sendRequest(
      address + "/api/v1/users/ignored",
      "POST",
      {
        "Content-Type": "application/json",
        Authorization: "Bearer " + accessToken,
      },
      {
        userId: userId,
        friendId: ignoredUser.id,
      }
    ).then((response) => {
      if (response.status.ok) {
        context.commit("addIgnoredUser", ignoredUser);
        context.commit("removeFriend", ignoredUser.id);
      } else {
        console.log(response.data.message);
      }
    });
  },
  removeFriend(context, friendId) {
    let userId = context.getters["userId"];
    let accessToken = context.getters["accessToken"];
    sendRequest(
      address + "/api/v1/users/friends",
      "DELETE",
      {
        "Content-Type": "application/json",
        Authorization: "Bearer " + accessToken,
      },
      {
        userId: userId,
        friendId: friendId,
      }
    ).then((response) => {
      if (response.status.ok) {
        context.commit("removeFriend", friendId);
      } else {
        console.log(response.data.message);
      }
    });
  },
  removeIgnored(context, ignoredUserId) {
    let userId = context.getters["userId"];
    let accessToken = context.getters["accessToken"];
    sendRequest(
      address + "/api/v1/users/ignored",
      "DELETE",
      {
        "Content-Type": "application/json",
        Authorization: "Bearer " + accessToken,
      },
      {
        userId: userId,
        friendId: ignoredUserId,
      }
    ).then((response) => {
      if (response.status.ok) {
        context.commit("removeIgnoredUser", ignoredUserId);
      } else {
        console.log(response.data.message);
      }
    });
  },
  changePassword(context, data) {
    return new Promise((resolve, reject) => {
      let userId = context.getters["userId"];
      let username = context.getters["username"];
      let accessToken = context.getters["accessToken"];
      sendRequest(
        address + "/api/v1/users/update/password",
        "PUT",
        {
          "Content-Type": "application/json",
          Authorization: "Bearer " + accessToken,
        },
        {
          username,
          userId,
          oldPassword: data.password,
          newPassword: data.newPassword,
        }
      ).then((response) => {
        if (response.status.ok) {
          context.commit("saveUserSession", response.data);
          resolve();
        } else {
          reject(response.data.message);
        }
      });
    });
  },
  changeUsername(context, data) {
    return new Promise((resolve, reject) => {
      let userId = context.getters["userId"];
      let username = context.getters["username"];
      let accessToken = context.getters["accessToken"];
      sendRequest(
        address + "/api/v1/users/update/username",
        "PUT",
        {
          "Content-Type": "application/json",
          Authorization: "Bearer " + accessToken,
        },
        {
          username,
          userId,
          password: data.password,
          newUsername: data.newUsername,
        }
      ).then((response) => {
        if (response.status.ok) {
          context.commit("updateUsername", data.newUsername);
          resolve();
        } else {
          reject(response.data.message);
        }
      });
    });
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

function saveToLocalStorage(data) {
  for (let field in data) {
    if (Array.isArray(data[field])) {
      window.localStorage.setItem(field, JSON.stringify(data[field]));
    } else {
      window.localStorage.setItem(field, data[field]);
    }
  }
}
