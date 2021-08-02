<template>
  <div>
    <h2>Change Password</h2>
    <form @submit.prevent="userChangePassword">
      <div class="formField">
        <label for="password">Current Password<abbr title="required">*</abbr>: </label>
        <input type="password" name="password" v-model.trim="password"  :class="{invalid: passwordInvalid}" @blur="validatePassword">
        <ul v-if="passwordInvalid" class="errorBox">
          <li v-for="e, i in passwordErrors" :key="i">{{e}}</li>
        </ul>
      </div>
      <div class="formField">
        <label for="newPassword">New Password<abbr title="required">*</abbr>: </label>
        <input type="password" name="newPassword" v-model.trim="newPassword"  :class="{invalid: newPasswordInvalid}" @blur="validateNewPassword">
        <ul v-if="newPasswordInvalid" class="errorBox">
          <li v-for="e, i in newPasswordErrors" :key="i">{{e}}</li>
        </ul>
      </div>
      <div class="formField">
        <label for="confirmPassword">Confirm New Password<abbr title="required">*</abbr>: </label>
        <input type="password" name="confirmPassword" v-model.trim="confirmPassword"  :class="{invalid: confirmPasswordInvalid}" @blur="validateConfirmPassword">
        <ul v-if="confirmPasswordInvalid" class="errorBox">
          <li v-for="e, i in confirmPasswordErrors" :key="i">{{e}}</li>
        </ul>
      </div>
        <button>Submit</button>
    </form>
  </div>
</template>

<script>
export default {
  data() {
    return {
      password: "",
      passwordInvalid: false,
      passwordErrors: [],
      newPassword: "",
      newPasswordInvalid: false,
      newPasswordErrors: [],
      confirmPassword: "",
      confirmPasswordInvalid: false,
      confirmPasswordErrors: [],
    }
  },
  methods: {
    userChangePassword() {
    // In case the user hits enter when in a field.
    this.validatePassword();
    this.validateNewPassword();
    this.validateConfirmPassword();

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
    },
    validatePassword() {
      this.passwordInvalid= false;
      this.passwordErrors = [];
      
      if (this.password === "") {
        this.passwordErrors.push("must not be empty");
        this.passwordInvalid = true;

        return
      } else {
        if ( !(7 < this.password.length && this.password.length < 25) ) {
          this.passwordErrors.push("should be between 8 and 24 characters");
          this.passwordInvalid = true;
        }
        if ( !/^(?=.*\d)(?=.*[A-Z])(?=.*[a-z])(?=.*\W|_)\S{1,}$/.test(this.password)) {
          console.log(this.password)
          this.passwordErrors.push("must have at least one of each:")
          this.passwordErrors.push("  - lowercase character")
          this.passwordErrors.push("  - uppercase character")
          this.passwordErrors.push("  - digit (0-9)")
          this.passwordErrors.push("  - special symbol: `~!@#$%^&*()_-+=[]{}\\|/:\";'<>,.")
          this.passwordInvalid = true;
        }
      }
    },
    validateNewPassword() {
      this.newPasswordInvalid= false;
      this.newPasswordErrors = [];
      
      if (this.newPassword === "") {
        this.newPasswordErrors.push("must not be empty");
        this.newPasswordInvalid = true;

        return
      } else {
        if ( !(7 < this.password.length && this.password.length < 25) ) {
          this.newPasswordErrors.push("should be between 8 and 24 characters");
          this.newPasswordInvalid = true;
        }
        if ( !/^(?=.*\d)(?=.*[A-Z])(?=.*[a-z])(?=.*\W|_)\S{1,}$/.test(this.password)) {
          this.newPasswordErrors.push("must have at least one of each:")
          this.newPasswordErrors.push("  - lowercase character")
          this.newPasswordErrors.push("  - uppercase character")
          this.newPasswordErrors.push("  - digit (0-9)")
          this.newPasswordErrors.push("  - special symbol: `~!@#$%^&*()_-+=[]{}\\|/:\";'<>,.")
          this.newPasswordInvalid = true;
        }
      }
    },
    validateConfirmPassword() {
      this.confirmPasswordInvalid = false;
      this.confirmPasswordErrors = [];
      
      if (this.confirmPassword !== this.newPassword) {
        this.confirmPasswordErrors.push("the passwords do not match");
        this.confirmPasswordInvalid = true;

        return
      }
    }
  }
}
</script>

<style scoped>
.invalid {
  border: 2px solid red;
}
.errorBox {
  list-style-type: none;
  margin: 2px 0 4px;
  padding: 0;
  color: red;
}
.errorBox > li {
  margin-left: 16px;
  padding: 2px 0;
}
</style>