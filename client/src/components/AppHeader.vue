<template>
  <div id="header">
    <p id="username">
      <router-link to="/">HOME</router-link>
      <span v-if="loggedIn">Username: {{ username }} | <router-link to="" @click.prevent="onLogout">Logout</router-link></span>
    </p>
    <ul class="app-header-navbar">
      <li class="app-header-navbar-item" v-if="!loggedIn">
        <router-link to="/signup">Create Account</router-link>
      </li>
      <li class="app-header-navbar-item" v-if="!loggedIn">
        <router-link to="/signin">Sign In</router-link>
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
    }
  },
  methods: {
    onLogout() {
      this.$store.dispatch("user/logout");
      this.$router.push("/");
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
</style>