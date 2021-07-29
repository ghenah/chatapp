<template>
  <div>
    <ul>
      <li v-for="n in notifications" :key="n.id">
        <slot name="notifItem"  :n="n"></slot>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  computed: {
    notifications() {
      console.log(this.$store.getters["chat/notifications"]);
      return this.$store.getters["chat/notifications"];
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