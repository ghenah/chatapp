<template>
  <div>
    <h3>Get User</h3>
    <form @submit.prevent="onSubmit">
      UserID
      <input id="get-user-id" name="get-user-id" v-model="id"/>
      <button>Get User</button>
    </form>
    <div>
      <ul>
        <li>User ID: {{ userId }}</li>
        <li>Username: {{ username }}</li>
        <li>Email: {{ email }}</li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  props: [
    'userId',
    'username',
    'email'
  ],
    data() {
      return {
        id: ""
      }
    },
  methods: {

    onSubmit() {
      this.$emit("get-user-by-id", this.id)
    }
  }
}
</script>