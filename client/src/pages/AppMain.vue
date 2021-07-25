<template>
  <div v-if="loggedIn">
    <main-chat class="page"></main-chat>
    <friends-list class="page"></friends-list>
    <ignore-list class="page"></ignore-list>
  </div>
</template>

<script>
import MainChat from "../components/chat/MainChat.vue";
import FriendsList from "../components/user-account/FriendsList.vue";
import IgnoreList from "../components/user-account/IgnoreList.vue";
export default {
  components: {
    MainChat,
    FriendsList,
    IgnoreList
  },
  computed: {
    loggedIn() {
      return this.$store.getters["user/loggedIn"];
    }
  }
}
</script>