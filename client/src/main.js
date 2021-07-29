import { createApp } from "vue";
import App from "./App";
const app = createApp(App);

import router from "./router";
import store from "./store/index";
import BaseContextMenu from "./components/base/BaseContextMenu.vue";
import BasePopupWindow from "./components/base/BasePopupWindow.vue";
import BaseButton from "./components/base/BaseButton.vue";

app.use(store);
app.use(router);

app.component("base-context-menu", BaseContextMenu);
app.component("base-popup-window", BasePopupWindow);
app.component("base-button", BaseButton);

app.mount("#app");
