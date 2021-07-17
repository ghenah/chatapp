import { createApp } from "vue";
import App from "./App";
const app = createApp(App);

import router from "./router";
import store from "./store/index";

app.use(store);
app.use(router);

app.mount("#app");
