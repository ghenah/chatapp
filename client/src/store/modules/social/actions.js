var schema = process.env.VUE_APP_SCHEMA;
var domain = process.env.VUE_APP_DOMAIN;
var port = process.env.VUE_APP_PORT;
var address = schema + domain + ":" + port;

export default {
  searchUsers(context) {
    let accessToken = context.rootGetters["user/accessToken"];
    sendRequest(address + "/api/v1/users/search", "GET", {
      "Content-Type": "application/json",
      Authorization: "Bearer " + accessToken,
    }).then((response) => {
      if (response.status.ok) {
        context.commit("saveUserSearchResults", response.data.usersList);
      } else {
        console.log(response.data.message);
      }
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
