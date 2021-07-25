<template>
  <div id="mainChatContainer">
    {{connected ? "Connected" : ""}}
    <button v-if="!connected" @click="onChatServerConnect">Connect</button>
    <form id="formNewRoom" @submit.prevent="onNewChatRoom">
      <select v-model="newRoomVisibility">
        <option :value="roomVisPublic">Public</option>
        <option :value="roomVisPrivate">Private</option>
      </select>
      <input type="submit" value="New Room">
      <input v-model="newRoomName" type="text" placeholder="Room name mandatory">
    </form>
    <div id="chatArea">
      <div id="roomsList">
        <ul>
          <li v-for="room in rooms" :key="room.id">
            <div @click="onChatRoomSelected(room.id)" :class="{'room-name-active': room.id === activeRoomId}" class="room-name">{{room.name}}</div>
            <ul v-if="room.id === activeRoomId">
              <li v-for="username, userId in room.activeUsers" :key="userId" class="users-list-item">{{username}}</li>
            </ul>
          </li>
        </ul>
      </div>
      <div id="chatAreaMain">
        <div>
          <button @click="onLeaveChatRoom">Leave</button>
        </div>
        <div id="messages">
          <div v-for="m in messages" :key="m.type" class="message">
            <p v-if="m.type === 'userJoined'">{{m.username}} joined</p>
            <p v-else-if="m.type === 'userLeft'">{{m.username}} left</p>
            <p v-else-if="m.type === 'newChatMessage'">{{m.authorUsername}} : {{m.content}}</p>
          </div>
        </div>
        <input id="chatTextArea" type="text" v-model="newOutMessage" @keyup.enter="onNewOutMessage">
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      newRoomVisibility: 1,
      newRoomName: "",
      newOutMessage: "",
      activeChatId: 0,
    }
  },
  computed: {
    roomVisPublic() {
      return this.$store.getters["chat/roomVisPublic"];
    },
    roomVisPrivate() {
      return this.$store.getters["chat/roomVisPrivate"];
    },
    rooms() {
      return this.$store.getters["chat/roomsList"];
    },
    messages() {
      return this.$store.getters["chat/messagesLog"]
    },
    activeRoomId() {
      return this.$store.getters["chat/activeRoomId"];
    },
    connected() {
      return this.$store.getters["chat/connected"];
    }
  },
  methods: {
    onChatServerConnect() {
      this.$store.dispatch("chat/getTicket")
      .then(() => {
        this.$store.dispatch("chat/openWS")
      })
      .catch((errorMsg) => {
        console.log("Error: fetch ticket: " + errorMsg);
      });
    },
    onNewChatRoom() {
      if (this.newRoomName === "") {
        alert("Make sure to also give the room a name.");
        return
      }

      this.$store.dispatch("chat/newRoom", {visibility: this.newRoomVisibility, name: this.newRoomName})
      this.newRoomVisibility = 1;
      this.newRoomName = "";
    },
    onChatRoomSelected(roomId) {
      this.newOutMessage = "";
      this.$store.dispatch("chat/setActiveChatRoom", {roomId});
    },
    onNewOutMessage() {
      this.$store.dispatch("chat/sendMessage", {roomId: this.activeRoomId, content: this.newOutMessage});
      this.newOutMessage = "";
    },
    onLeaveChatRoom() {
      this.$store.dispatch("chat/leaveChatRoom", {roomId: this.activeRoomId});
    }
  }
}
</script>

<style scoped>
  #mainChatContainer {
    margin: 0;
    padding: 0;
  }
ul {
  list-style-type: none;
  margin: 0;
  padding: 0;
}
.room-name {
  padding-left: 8px;
}
.room-name-active {
  color: #f0f0f0;
  background-color: #202020;
}
.users-list-item {
  margin-left: 16px;
}
#chatArea {
  display: flex;
  min-height: 200px;
}
#roomsList {
  display: flexbox;
  width: 180px;
  border-right: solid 1px black;
}
#chatAreaMain {
  display: flex;
  flex-grow: 1;
  flex-direction: column;
}
#messages {
  display: flexbox;
  flex-grow: 1;
}
.message > p {
  margin: 0;
  padding: 4px;
}
#chatTextArea {
  display: flexbox;
  height: 10%;
  border-radius: 0px;
  border: 0px;
  border-top: 1px solid black
}
#formNewRoom {
  border: solid 1px black;
  border-left: 0px;
  border-right: 0px;
  padding: 8px;
}
#formNewRoom > * {
  margin-right: 4px;
}
</style>