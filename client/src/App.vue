<template>
  <div>
      <app-header></app-header>
      <main>
        <router-view></router-view>
      </main>
  </div>
</template>

<script>
import AppHeader from "./components/AppHeader.vue"

export default {
  components: {
    AppHeader
  },
  created() {
    if (this.$store.getters["user/loggedIn"]) {
      // The status must be set to logged out so that the app does not
      // start rendering the content.
      this.$store.dispatch("user/setStatusLoggedOut");

      // This means a user login session is stored in the local storage
      // and the user profile information must be fetched.
      this.$store.dispatch("user/getProfileInfo")
      .then(() => {
        // Now the app can render the content.
        this.$store.dispatch("user/setStatusLoggedIn");
      })
      .catch((errorMsg) => {
        console.log(errorMsg);
        this.$store.dispatch("user/logout");
      })
    }
  }

}
</script>

<style>
  button {
    margin-right: 4px;
  }
  div.formField {
    margin-top: 4px;
  }
  abbr {
    text-decoration: none;
  }
  div.page {
    border: 1px solid black;
    margin: 10px auto;
    padding: 6px 8px;
  }
  #userInfoField {
    height: 1rem;
    padding-top: 2px;
    padding-bottom: 2px;
  }
</style>