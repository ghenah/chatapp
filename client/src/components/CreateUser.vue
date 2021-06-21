<template>
  <form @submit.prevent="onSubmit">
    <h3>Create User</h3>
    <div>
      User ID
      <input id="id" name="id" type="text" v-model="id"/>
    </div>
    <div>
      Username
      <input id="username" name="username" type="text" v-model="username"/>
    </div>
    <div>
      Email
      <input id="email" name="email" type="text" v-model="email"/>
    </div>
    <button>Create User</button>
  </form>

</template>

<script>
export default {
  methods: {
    data() {
      return {
        id: "",
        username: "",
        email: ""
      }
    },
    onSubmit() {
      this.$emit("create-user-submitted", this.id, this.username, this.email)
    }
  }
}
</script>