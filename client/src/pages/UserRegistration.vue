<template>
  <div v-if="!loggedIn">
    <h2>Create an Account</h2>
    <form @submit.prevent="onFormSubmit">
      <div class="formField">
        <label for="username">Username<abbr title="required">*</abbr>: </label>
        <input type="text" name="username" v-model="username">
      </div>
      <div class="formField">
        <label for="email">Email<abbr title="required">*</abbr>: </label>
        <input type="text" name="email" v-model="email">
      </div>
      <div class="formField">
        <label for="password">Password<abbr title="required">*</abbr>: </label>
        <input type="password" name="password" v-model="password">
      </div>
      <div class="formField">
        <label for="repeatPassword">Repeat Password<abbr title="required">*</abbr>: </label>
        <input type="password" name="repeatPassword" v-model="repeatPassword">
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
      email: "",
      password: "",
      repeatPassword: ""
    }
  },
  methods: {
    onFormSubmit() {
      if (this.password !== this.repeatPassword) {
        console.log("Passwords don't match");
        return;
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
    }
  },
  beforeCreate() {
    if (this.$store.getters["loggedIn"]) {
      this.$router.push("/");
    }
  }
  
}
</script>