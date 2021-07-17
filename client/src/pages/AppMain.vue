<template>
  <div v-if="loggedIn">
    <friends-list class="page"></friends-list>
    <ignore-list class="page"></ignore-list>
  </div>
</template>

<script>
import FriendsList from "../components/user-account/FriendsList.vue";
import IgnoreList from "../components/user-account/IgnoreList.vue";
export default {
  components: {
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