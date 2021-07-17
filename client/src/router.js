import { createRouter, createWebHistory } from "vue-router";

import AppMain from "./pages/AppMain.vue";
import UserRegistration from "./pages/UserRegistration.vue";
import UserAuthentication from "./pages/UserAuthentication.vue";
import AccountSettings from "./pages/AccountSettings.vue";
import SearchUsers from "./pages/SearchUsers.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", component: AppMain },
    { path: "/signup", component: UserRegistration },
    { path: "/signin", component: UserAuthentication },
    { path: "/users/search", component: SearchUsers },
    { path: "/account/settings", component: AccountSettings },
  ],
});

export default router;
