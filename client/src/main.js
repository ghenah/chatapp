import { createApp } from "vue";
import App from "./App";
const app = createApp(App);

import router from "./router";
import store from "./store/index";
import BaseContextMenu from "./components/base/BaseContextMenu.vue";

app.use(store);
app.use(router);

app.component("base-context-menu", BaseContextMenu);

app.mount("#app");
