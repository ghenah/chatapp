<template>
  <div v-if="!loggedIn">
    <h2>Sign In</h2>
    <form @submit.prevent="onFormSubmit">
      <div class="formField">
        <label for="username">Username: </label>
        <input type="text" name="username" v-model="username">
      </div>
      <div class="formField">
        <label for="password">Password: </label>
        <input type="password" name="password" v-model="password">
      </div>
      <div>
        <button>Sign In</button>
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
      password: ""
    }
  },
  methods: {
    onFormSubmit() {
      this.$store.dispatch("user/login", {
        username: this.username,
        password: this.password
      })
      .then(() => {
        this.password = "";
        this.username = "";

        this.$router.push("/");
      })
      .catch((errorMsg) => {
        this.password = "";
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