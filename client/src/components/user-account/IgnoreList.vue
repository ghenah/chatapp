<template>
  <div>
    <h2>Ignored[{{ignoredList.length}}]</h2>
    <ul>
      <li v-for="ignoredUser in ignoredList" :key="ignoredUser.id">{{ignoredUser.username}} <button @click="removeIgnoredUser(ignoredUser.id)">X</button> </li>
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
  methods: {
    removeIgnoredUser(ignoredUserId) {
      this.$store.dispatch("user/removeIgnored", ignoredUserId);
    }
  }
}
</script>

<style scoped>
button {
  padding: 0px;
  height: 1rem;
  width: 1rem;
}
</style>