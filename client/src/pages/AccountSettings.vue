<template>
  <div v-if="loggedIn">
    <change-password></change-password>
    <change-username></change-username>
  </div>
</template>

<script>
import ChangePassword from "../components/user-account/ChangePassword.vue";
import ChangeUsername from "../components/user-account/ChangeUsername.vue";

export default {
  computed: {
    loggedIn() {
      return this.$store.getters["user/loggedIn"];
    }
  },
  components: {
    ChangePassword,
    ChangeUsername
  },
}
</script>