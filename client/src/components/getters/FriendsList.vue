<template>
  <div>
    <ul>
      <li v-for="friend in friends" :key="friend.id">
        <slot name="friendItem" :frUsername="friend.username" :frId="friend.id">{{friend.username}}</slot>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  computed: {
    friends() {
      return this.$store.getters["user/friendsList"];
    }
  }
}
</script>

<style scoped>
ul {
  list-style-type: none;
  margin: 0px;
  padding: 0px;
}
</style>