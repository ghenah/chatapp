var schema = process.env.VUE_APP_SCHEMA;
var domain = process.env.VUE_APP_DOMAIN;
var port = process.env.VUE_APP_PORT;
var address = schema + domain + ":" + port;

export default {
  submit(context, formData) {
    return new Promise((resolve, reject) => {
      sendRequest(
        address + "/auth/signup",
        "POST",
        { "Content-Type": "application/json" },
        {
          email: formData.email,
          username: formData.username,
          password: formData.password,
        }
      ).then((response) => {
        if (response.status.ok) {
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
