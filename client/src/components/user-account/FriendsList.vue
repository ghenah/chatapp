<template>
  <div>
    <h2>Friends[{{friends.length}}]</h2>
    <ul>
      <li v-for="friend in friends" :key="friend.id">{{friend.username}} <button @click="removeFriend(friend.id)">Remove</button></li>
    </ul>
  </div>
</template>

<script>
export default {
  computed: {
    friends() {
      return this.$store.getters["user/friendsList"];

    }
  },
  methods: {
    removeFriend(friendId) {
      this.$store.dispatch("user/removeFriend", friendId)
    }
  }
}
</script>