<template>
  <div id="header">
    <p id="username">
      <router-link to="/" class="home-link">HOME</router-link>
      <span v-if="loggedIn">{{ username }} | <router-link to="" @click.prevent="onLogout">sign out</router-link></span>
    </p>
    <ul class="app-header-navbar">
      <li class="app-header-navbar-item" v-if="!loggedIn">
        <router-link to="/signup">Create Account</router-link>
      </li>
      <li class="app-header-navbar-item" v-if="!loggedIn">
        <router-link to="/signin">Sign In</router-link>
      </li>
      <li class="app-header-navbar-item" v-if="loggedIn && !notifPermissionGranted">
        <div role="button" id="notif-enable" @click="requestPermission">Enable Notifications</div>
      </li>
      <li class="app-header-navbar-item" v-if="loggedIn">
        <router-link to="/account/settings">Account</router-link>
      </li>
      <li class="app-header-navbar-item" v-if="loggedIn">
        <router-link to="/users/search">Search Users</router-link>
      </li>
      <li class="app-header-navbar-item" v-if="loggedIn">
        <router-link to="/chat/search">Search Chat Rooms</router-link>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  computed: {
    username() {
      return this.$store.getters["user/username"];
    },
    loggedIn() {
      return this.$store.getters["user/loggedIn"];
    },
    notifPermissionGranted() {
      if (Notification.permission === "granted") {
        return true;
      } else {
        return false;
      }
      
    },
  },
  methods: {
    onLogout() {
      this.$store.dispatch("user/logout");
      this.$router.push("/");
    },
    requestPermission() {
      Notification.requestPermission().then((permission) => {
        if (permission === "granted") {
          const n = new Notification("Hello!", {
            body: "This is a message body."
          });

          n.close();
        }
      })
    }
  }
}
</script>

<style>
#header {
  display: flex;
}
.app-header-navbar {
  display: flex;
  list-style: none;
  margin-left: auto;
}
.app-header-navbar-item {
  padding: 2px 6px;
}
.app-header-navbar-item:hover {
  background-color: #e8e8e8;
}
#notif-enable {
  display: inline-block;
  font-size: 12px;
  font-family: sans-serif;
  background: hsl(18, 80%, 50%);
  padding: 4px;
  color: white;
}
.home-link {
  margin-right: 14px;
}
</style>