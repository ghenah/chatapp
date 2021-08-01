<template>
  <div id="mainChatContainer">
    <div>
      {{connected ? "Connected" : ""}}
      <div role="button" v-if="!connected" @click="onChatServerConnect" id="connectBtn">Connect</div>
    </div>
    <div id="chatControls">
      <form id="formNewRoom" @submit.prevent="onNewChatRoom">
        <select v-model="newRoomVisibility">
          <option :value="roomVisPublic">Public</option>
          <option :value="roomVisPrivate">Private</option>
        </select>
        <input type="submit" value="New Room">
        <input v-model="newRoomName" type="text" placeholder="Room name mandatory">
      </form>
      <div v-if="notifNumber > 0" id="notificationsIcon" @click="onNotifIconClick">{{notifNumber}}</div>
    </div>
    <div id="chatArea">
      <div id="roomsList">
        <ul>
          <li v-for="room in rooms" :key="room.id">
            <div @click="onChatRoomSelected(room.id)" @contextmenu.prevent="onChatRoomCM($event, room.id, room.name)" :class="{'room-name-active': room.id === activeRoomId}" class="room-name">{{room.name}}</div>
            <ul v-if="room.id === activeRoomId">
              <!-- <li v-for="username, userId in room.activeUsers" :key="userId" class="users-list-item">{{username}}</li> -->
              <li v-for="username, userId in room.activeUsers" :key="userId" class="users-list-item">
                <img :src="profilePics.get(userId)" alt="Profile picture" class="profile-pic-micro">
                {{username}}</li>
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
    <base-context-menu ref="contextMenu" :cfg="cmCfg" :actions="cmActions" @closed="cmCleanup">
      <template #menuName="{cmName}">
        {{cmName}}
      </template>
      <template #action="{actionName, cb}">
        <div @click="cb">{{actionName}}</div>
      </template>
    </base-context-menu>
    <base-popup-window ref="popupWindow" :dimensions="pwUsersList.dimensions">
      <div class="pwListSwitcher">
        <div class="pwListSwitcherItem" @click="onInviteListSwitch('search')">Search</div>
        <div class="pwListSwitcherItem" @click="onInviteListSwitch('friends')">Friends</div>
      </div>
      <users-search v-if="pwCurrentInviteList === 'search' ">
        <template #userItem="{uUsername, uId}">
          <div class="plUserItem">
            <div class="plUserItemName">
              {{uUsername}}
            </div>
            <div class="plUserItemActions">
              <div @click="onUserInvite(uId, uUsername)">Invite</div>
            </div>
          </div>
        </template>
      </users-search>
      <friends-list v-if="pwCurrentInviteList === 'friends' ">
        <template #friendItem="{frUsername, frId}">
          <div class="plUserItem">
            <div class="plUserItemName">
              {{frUsername}}
            </div>
            <div class="plUserItemActions">
              <div @click="onUserInvite(frId, frUsername)">Invite</div>
            </div>
          </div>
        </template>
      </friends-list>
    </base-popup-window>
    <base-popup-window ref="pwNotifications" :dimensions="pwNotifications.dimensions">
      <notifications>
        <template #notifItem="{n}">
          <div class="pwNotifItem">
            <div class="pwInvInfo">
              <div class="pwInvInfoHeader">{{n.username}} invites you to join:</div>
              <div class="pwInvInfoBody">{{n.roomName}}</div>
            </div>
            <div class="pwInvInfoAcceptDecline">
              <base-button :btnType="'default'" :btnClass="'positive'">
                <template #default>
                  <div @click="onAcceptInvitation(n.id, n.roomId, true)">Accept</div>
                </template>
              </base-button>
              <base-button :btnType="'default'" :btnClass="'negative'">
                <template #default>
                  <div @click="onAcceptInvitation(n.id, n.roomId, false)">Decline</div>
                </template>
              </base-button>
            </div>
          </div>
        </template>
      </notifications>
    </base-popup-window>
  </div>
</template>

<script>
import UsersSearch from "../getters/UsersSearch.vue";
import FriendsList from "../getters/FriendsList.vue";
import Notifications from "../getters/Notifications.vue";

export default {
  components: {
    UsersSearch,
    FriendsList,
    Notifications
  },
  data() {
    return {
      newRoomVisibility: 1,
      newRoomName: "",
      newOutMessage: "",
      activeChatId: 0,
      cmCfg: {
        hadName: false,
        name: "",
      },
      cmActions: [],
      pwUsersList: {
        roomId: 0,
        roomName: "",
        dimensions: {
          width: 240,
          height: 380
        }
      },
      pwNotifications: {
        dimensions: {
          width: 220,
          height: 300,
        }
      },
      pwCurrentInviteList: "friends",
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
    profilePics() {
      return this.$store.getters["chat/userProfilePics"];
    },
    activeRoomId() {
      return this.$store.getters["chat/activeRoomId"];
    },
    connected() {
      return this.$store.getters["chat/connected"];
    },
    notifNumber() {
      return this.$store.getters["chat/notifications"].length;
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
    },
    onUserInvite(inviteeId, inviteeUsername) {
      this.$store.dispatch("chat/inviteUser", {inviteeId, inviteeUsername, roomId: this.pwUsersList.roomId});
    },
    onAcceptInvitation(invId, roomId, accepted) {
      this.$store.dispatch("chat/acceptInvitation", {invId, roomId, accepted})
    },
    onInviteUsersWindow(roomId, roomName) {
      this.pwUsersList.roomId = roomId;
      this.pwUsersList.roomName = roomName;
      this.$refs.popupWindow.open();
    },
    onChatRoomCM(e, roomId, roomName) {
      this.cmCfg = {
        hasName: true,
        name: roomName
      }
      this.cmActions = [
        {
          name: "Invite users",
          action: () => {
            this.$refs.contextMenu.close();
            this.onInviteUsersWindow(roomId);
          }
        }
      ]

      this.$refs.contextMenu.open(e)
    },
    cmCleanup() {
      this.cmCfg = {};
      this.cmActions = [];
    },
    onInviteListSwitch(listName) {
      this.pwCurrentInviteList = listName;
    },
    onNotifIconClick() {
      this.$refs.pwNotifications.open();
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
  padding: 8px;
}
#formNewRoom > * {
  margin-right: 4px;
}
.plUserItem {
  padding: 4px;
  display: flex;
}
.plUserItem:hover {
  background-color: hsl(230, 30%, 70%);
}
.plUserItem .plUserItemActions {
  margin-left: auto;
  visibility: hidden;
}
.plUserItem:hover .plUserItemActions {
  visibility: visible;
}
.pwListSwitcher {
  display: flex;
  justify-content: space-evenly;
}
.pwListSwitcherItem {
  padding: 4px;
  display: flexbox;
  text-decoration: underline;
}
.pwNotifItem {
  display: flex
}
.pwInvInfo {
  width: 160px
}
.pwInvInfoHeader {
  text-align: center;
  font-size: 12px;
}
.pwInvInfoBody {
  text-align: center;
  font-size: 14px;
}
.pwInvInfoAcceptDecline {
  width: 60px;
  display: flex;
  flex-direction: column;
}
#chatControls {
  display: flex;
  border: solid 1px black;
  border-left: 0px;
  border-right: 0px;
}
#notificationsIcon {
  text-align: center;
  width: 20px;
  margin-left: auto;
  background: #ff3333;
  color: white;
  font-size: 18px;
  font-weight: bold;
}
#connectBtn {
  display: inline-block;
  padding: 2px 4px;
  font-family: sans-serif;
  height: 100%;
  border-width: 0px;
  background: #ff3333;
  color: white;
}
.profile-pic-micro {
  height: 50px;
  width: 50px;
  border-radius: 25px;
}
</style>