<template>
  <div>
    <ul>
      <li v-for="ignoredUser in ignoredList" :key="ignoredUser.id">
        <slot name="ignoredItem" :iUsername="ignoredUser.username" :iId="ignoredUser.id">{{ignoredUser.username}}</slot>
        </li>
    </ul>
  </div>
</template>

<script>
export default {
  computed: {
    ignoredList() {
      return this.$store.getters["user/ignoreList"];
    }
  },
}
</script>

<style scoped>
ul {
  list-style-type: none;
  margin: 0px;
  padding: 0px;
}
</style>