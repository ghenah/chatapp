<template>
  <div v-if="!loggedIn">
    <h2>Create an Account</h2>
    <form @submit.prevent="onFormSubmit">
      <div class="formField">
        <label for="username">Username<abbr title="required">*</abbr>: </label>
        <input type="text" name="username" v-model.trim="username" :class="{invalid: usernameInvalid}" @blur="validateUsername">
        <ul v-if="usernameInvalid" class="errorBox">
          <li v-for="e, i in usernameErrors" :key="i">{{e}}</li>
        </ul>
      </div>
      <div class="formField">
        <label for="email">Email<abbr title="required">*</abbr>: </label>
        <input type="text" name="email" v-model.trim="email" :class="{invalid: emailInvalid}" @blur="validateEmail">
        <ul v-if="emailInvalid" class="errorBox">
          <li v-for="e, i in emailErrors" :key="i">{{e}}</li>
        </ul>
      </div>
      <div class="formField">
        <label for="password">Password<abbr title="required">*</abbr>: </label>
        <input type="password" name="password" v-model.trim="password"  :class="{invalid: passwordInvalid}" @blur="validatePassword">
        <ul v-if="passwordInvalid" class="errorBox">
          <li v-for="e, i in passwordErrors" :key="i">{{e}}</li>
        </ul>
      </div>
      <div class="formField">
        <label for="repeatPassword">Repeat Password<abbr title="required">*</abbr>: </label>
        <input type="password" name="repeatPassword" v-model.trim="repeatPassword"  :class="{invalid: repeatPasswordInvalid}" @blur="validateRepeatPassword">
        <ul v-if="repeatPasswordInvalid" class="errorBox">
          <li v-for="e, i in repeatPasswordErrors" :key="i">{{e}}</li>
        </ul>
      </div>
      <div class="formField">
        <button>Create</button>
      </div>
    </form>
  </div>

</template>

<script>
export default {
  computed: {
    loggedIn() {
      return this.$store.getters["user/loggedIn"];
    }
  },
  data() {
    return {
      username: "",
      usernameInvalid: false,
      usernameErrors: [],
      email: "",
      emailInvalid: false,
      emailErrors: [],
      password: "",
      passwordInvalid: false,
      passwordErrors: [],
      repeatPassword: "",
      repeatPasswordInvalid: false,
      repeatPasswordErrors: [],
    }
  },
  methods: {
    onFormSubmit() {
    // In case the user hits enter when in a field.
    this.validateUsername();
    this.validateEmail();
    this.validatePassword();
    this.validateRepeatPassword();

      if (this.usernameInvalid || this.emailInvalid || this.passwordInvalid || this.repeatPasswordInvalid) {
        return
      }

      this.$store.dispatch("registration/submit", {
        username: this.username,
        email: this.email,
        password: this.password,
        repeatPassword: this.repeatPassword
      })
      .then(() => {
        this.username = "";
        this.email = "";
        this.password = "";
        this.repeatPassword = "";

        this.$router.push("/signin");
      })
      .catch((errorMsg) => {
        this.password = "";
        this.repeatPassword = "";
        console.log("Error: sign in:" + errorMsg);
      });
    },
    validateUsername() {
      this.usernameInvalid = false;
      this.usernameErrors = [];

      if (this.username === "") {
        this.usernameErrors.push("must not be empty");
        this.usernameInvalid = true;

        return
      } else {
        if ( !(1 < this.username.length && this.username.length < 17) ) {
          this.usernameErrors.push("should be between 2 and 16 characters")
          this.usernameInvalid = true;

          return
        }
        if ( !/^[A-Za-z]{1}[A-Za-z0-9]{1,15}$/.test(this.username)) {
          this.usernameErrors.push("can only contain letters and digits and must start with a letter")
          this.usernameInvalid = true;
        }
      }
    },
    validateEmail() {
      this.emailInvalid = false;
      this.emailErrors = [];
      
      if (this.email === "") {
        this.emailErrors.push("must not be empty");
        this.emailInvalid = true;

        return
      } else {
        if ( !/^[^\s@]+@[^\s@]+$/.test(this.email)) {
          this.emailErrors.push("should be in the form name@host")
          this.emailInvalid = true;
        }
      }
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
    validateRepeatPassword() {
      this.repeatPasswordInvalid = false;
      this.repeatPasswordErrors = [];
      
      if (this.repeatPassword !== this.password) {
        this.repeatPasswordErrors.push("the passwords do not match");
        this.repeatPasswordInvalid = true;

        return
      }
    }
  },
  beforeCreate() {
    if (this.$store.getters["loggedIn"]) {
      this.$router.push("/");
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