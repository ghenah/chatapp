import { createStore } from "vuex";

import userModule from "./modules/user/index.js";
import registrationModule from "./modules/registration/index.js";
import socialModule from "./modules/social/index.js";
import chatModule from "./modules/chat/index.js";
import chatSocket from "./plugins/chat.js";

const store = createStore({
  modules: {
    user: userModule,
    registration: registrationModule,
    social: socialModule,
    chat: chatModule,
  },
  plugins: [chatSocket()],
});

export default store;
