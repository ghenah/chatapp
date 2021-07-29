<template>
  <div v-if="loggedIn">
    <main-chat class="page"></main-chat>
    <friends-list class="page">
      <template #friendItem="{frUsername, frId}">
        <div @contextmenu.prevent="onFriendCM($event, frId, frUsername)">{{frUsername}}</div>
      </template>
    </friends-list>
    <ignore-list class="page">
      <template #ignoredItem="{iUsername, iId}">
        <div @contextmenu.prevent="onIgnoreCM($event, iId, iUsername)">{{iUsername}}</div>
      </template>
    </ignore-list>
    <base-context-menu ref="contextMenu" :cfg="elementCfg" :actions="elementActions" @closed="cmCleanup">
      <template #menuName="{cmName}">
        {{cmName}}
      </template>
      <template #action="{actionName, cb}">
        <div @click="cb">{{actionName}}</div>
      </template>
    </base-context-menu>
  </div>
</template>

<script>
import MainChat from "../components/chat/MainChat.vue";
import FriendsList from "../components/getters/FriendsList.vue";
import IgnoreList from "../components/getters/IgnoreList.vue";
export default {
  components: {
    MainChat,
    FriendsList,
    IgnoreList
  },
  data() {
    return {
      elementCfg: {
        hasName: false,
        name: "",
      },
      elementActions: [],
    }
  },
  computed: {
    loggedIn() {
      return this.$store.getters["user/loggedIn"];
    }
  },
  methods: {
    onFriendCM(e, userId, username) {
      this.elementCfg = {
        hasName: true,
        name: username
      }
      this.elementActions = [
        {
          name:"Personal Chat",
          action: () => {
            this.$refs.contextMenu.close();
            this.$store.dispatch("chat/startPersonalChat", {userId, username});
          }
        },
        {
          name:"Remove",
          action: () => {
            this.$refs.contextMenu.close();
            this.$store.dispatch("user/removeFriend", userId);
          }
        }
      ]

      this.$refs.contextMenu.open(e);
    },
      onIgnoreCM(e, userId, username) {
      this.elementCfg = {
        hasName: true,
        name: username
      }
      this.elementActions = [
        {
          name:"Remove",
          action: () => {
            this.$refs.contextMenu.close();
            this.$store.dispatch("user/removeIgnored", userId);
          }
        }
      ]

      this.$refs.contextMenu.open(e);
    },
    cmCleanup() {
      this.elementCfg = {};
      this.elementActions = [];
    }
  }
}
</script>

<style>
body {
  background-color: #d0d0d0;
}
</style>