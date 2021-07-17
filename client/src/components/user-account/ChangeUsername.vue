<template>
  <div>
    <h2>Change Username</h2>
    <form @submit.prevent="userChangeUsername">
      <div class="formField">
        <label for="password">Password<abbr title="required">*</abbr>: </label>
        <input type="password" name="password" v-model="password">
      </div>
      <div class="formField">
        <label for="new-username">New Username<abbr title="required">*</abbr>: </label>
        <input type="text" name="new-username" v-model="newUsername">
      </div>
      <div>
        <button>Submit</button>
      </div>
    </form>
  </div>
</template>

<script>
export default {
  data() {
    return {
      password: "",
      newUsername: ""
    }
  },
  methods: {
    userChangeUsername() {
      this.$store.dispatch("user/changeUsername", {password: this.password, newUsername: this.newUsername} )
      .then(() => {
        this.password = "";
        this.newUsername = "";
        console.log("Username changed successfully");
      })
      .catch((errorMsg) => {
        console.log("Error: change username: " + errorMsg);
        this.password = "";
      });
    }
  }
}
</script>