<template>
  <div>
    <h2>Ignored[{{ignoredList.length}}]</h2>
    <ul>
      <li v-for="ignoredUser in ignoredList" :key="ignoredUser.id">{{ignoredUser.username}} <button @click="removeIgnoredUser(ignoredUser.id)">Remove</button> </li>
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