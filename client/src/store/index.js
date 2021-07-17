import { createStore } from "vuex";

import userModule from "./modules/user/index.js";
import registrationModule from "./modules/registration/index.js";
import socialModule from "./modules/social/index.js";

const store = createStore({
  modules: {
    user: userModule,
    registration: registrationModule,
    social: socialModule,
  },
});

export default store;
