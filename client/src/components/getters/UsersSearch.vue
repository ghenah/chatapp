<template>
  <div>
    <base-button btnType="default" btnClass="default">
      <template #default>
        <div @click="searchUsers">Search</div>
      </template>
    </base-button>
    <ul>
      <li v-for="user in usersList" :key="user.id">
        <slot name="userItem" :uUsername="user.username" :uId="user.id">id:{{user.id}} | username:{{user.username}}</slot>
      </li>
    </ul>
  </div>
</template>

<script>
import BaseButton from "../base/BaseButton.vue"

export default {
  components: {
    BaseButton
  },
  data() {
    return {
      red: "red"
    }
  },
  computed: {
    usersList() {
      return this.$store.getters["social/userSearchResults"]
    },
  },
  methods: {
    searchUsers() {
      this.$store.dispatch("social/searchUsers");
    },
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