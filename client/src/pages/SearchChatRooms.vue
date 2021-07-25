<template>
  <div v-if="loggedIn">
    <h2>Search Chat Rooms</h2>
    <button @click="onSearchChatRooms">Search</button>
    <ul>
      <li v-for="room in chatRoomsList" :key="room.id">{{room.ownerUsername}} : {{room.name}}<button @click="onJoinChatRoom(room.id)">Join</button></li>
    </ul>
  </div>
</template>

<script>
export default {
  computed: {
    chatRoomsList() {
      return this.$store.getters["chat/chatRoomSearchResults"]
    },
    loggedIn() {
      return this.$store.getters["user/loggedIn"];
    }
  },
  methods: {
    onSearchChatRooms() {
      this.$store.dispatch("chat/searchChatRooms");
    },
    onJoinChatRoom(roomId) {
      this.$store.dispatch("chat/joinChatRoom", {roomId})
    }
  }
}
</script>