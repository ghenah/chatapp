<template>
  <div>
    <h2>Change Username</h2>
    <form @submit.prevent="userChangeUsername">
      <div class="formField">
        <label for="password">Password<abbr title="required">*</abbr>: </label>
        <input type="password" name="password" v-model.trim="password"  :class="{invalid: passwordInvalid}" @blur="validatePassword">
        <ul v-if="passwordInvalid" class="errorBox">
          <li v-for="e, i in passwordErrors" :key="i">{{e}}</li>
        </ul>
      </div>
      <div class="formField">
        <label for="newUsername">New Username<abbr title="required">*</abbr>: </label>
        <input type="text" name="newUsername" v-model.trim="newUsername" :class="{invalid: newUsernameInvalid}" @blur="validateNewUsername">
        <ul v-if="newUsernameInvalid" class="errorBox">
          <li v-for="e, i in newUsernameErrors" :key="i">{{e}}</li>
        </ul>
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
      passwordInvalid: false,
      passwordErrors: [],
      newUsername: "",
      newUsernameInvalid: false,
      newUsernameErrors: [],
    }
  },
  methods: {
    userChangeUsername() {

    this.validatePassword();
    this.validateNewUsername();


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
          this.passwordErrors.push("must have at least one of each:")
          this.passwordErrors.push("  - lowercase character")
          this.passwordErrors.push("  - uppercase character")
          this.passwordErrors.push("  - digit (0-9)")
          this.passwordErrors.push("  - special symbol: `~!@#$%^&*()_-+=[]{}\\|/:\";'<>,.")
          this.passwordInvalid = true;
        }
      }
    },
    validateNewUsername() {
      this.newUsernameInvalid = false;
      this.newUsernameErrors = [];

      if (this.newUsername === "") {
        this.newUsernameErrors.push("must not be empty");
        this.newUsernameInvalid = true;

        return
      } else {
        if ( !(1 < this.newUsername.length && this.newUsername.length < 17) ) {
          this.newUsernameErrors.push("should be between 2 and 16 characters")
          this.newUsernameInvalid = true;

          return
        }
        if ( !/^[A-Za-z]{1}[A-Za-z0-9]{1,15}$/.test(this.newUsername)) {
          this.newUsernameErrors.push("can only contain letters and digits and must start with a letter")
          this.newUsernameInvalid = true;
        }
      }
    },
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