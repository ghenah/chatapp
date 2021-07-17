<template>
  <div v-if="loggedIn">
    <h2>Search Users</h2>
    <button @click="searchUsers">Search</button>
    <ul>
      <li v-for="user in usersList" :key="user.id">{{user.id}} : {{user.username}} <button @click="addFriend(user.id, user.username)">Add Friend</button> <button @click="ignoreUser(user.id, user.username)">Ignore</button> </li>
    </ul>
  </div>
</template>

<script>
export default {
  computed: {
    usersList() {
      return this.$store.getters["social/userSearchResults"]
    },
    loggedIn() {
      return this.$store.getters["user/loggedIn"];
    }
  },
  methods: {
    searchUsers() {
      this.$store.dispatch("social/searchUsers");
    },
    addFriend(id, username) {
      this.$store.dispatch("user/addFriend", {id, username});
    },
    ignoreUser(id, username) {
      this.$store.dispatch("user/ignoreUser", {id, username});
    }
  }
}
</script>