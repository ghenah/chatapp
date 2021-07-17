<template>
  <div>
    <h2>Change Password</h2>
    <form @submit.prevent="userChangePassword">
      <div class="formField">
        <label for="password">Current Password<abbr title="required">*</abbr>: </label>
        <input type="password" name="password" v-model="password">
      </div>
      <div class="formField">
        <label for="new-password">New Password<abbr title="required">*</abbr>: </label>
        <input type="password" name="new-password" v-model="newPassword">
      </div>
      <div class="formField">
        <label for="confirm-password">Confirm New Password<abbr title="required">*</abbr>: </label>
        <input type="password" name="confirm-password" v-model="confirmPassword">
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
      newPassword: "",
      confirmPassword: ""
    }
  },
  methods: {
    userChangePassword() {
      if (this.newPassword !== this.confirmPassword) {
        console.log("Passwords don't match");
        return;
      }

      this.$store.dispatch("user/changePassword", {password: this.password, newPassword: this.newPassword})
      .then(() => {
        console.log("Password updated successfully");
      })
      .catch((errorMsg) => {
        console.log("Error: change password: " + errorMsg);
      });
      this.password = "";
      this.newPassword = "";
      this.confirmPassword = "";      
    }
  }
}
</script>