<template>
  <div v-if="loggedIn">
    <base-context-menu ref="contextMenu" :cfg="elementCfg" :actions="elementActions" @closed="cmCleanup">
      <template #menuName="{cmName}">
        {{cmName}}
      </template>
      <template #action="{actionName, cb}">
        <div @click="cb">{{actionName}}</div>
      </template>
    </base-context-menu>
    <h2>Search Users</h2>
    <users-search>
      <template #userItem="{uUsername, uId}">
        <div class="usersListItem" @contextmenu.prevent="onUserCM($event, uId, uUsername)">{{uUsername}}</div>
      </template>
    </users-search>
  </div>
</template>

<script>
import UsersSearch from "../components/getters/UsersSearch.vue";

export default {
  components: {
    UsersSearch
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
    addFriend(id, username) {
      this.$store.dispatch("user/addFriend", {id, username});
    },
    ignoreUser(id, username) {
      this.$store.dispatch("user/ignoreUser", {id, username});
    },
    onUserCM(e, userId, username) {
      this.elementCfg = {
        hasName: true,
        name: username
      }
      this.elementActions = [
        { name:"Add as Friend",
          action: () => { 
            this.$refs.contextMenu.close();
            this.addFriend(userId, username);
          }  
        },
        { name:"Ignore User",
          action: () => { 
            this.$refs.contextMenu.close();
            this.ignoreUser(userId, username);
          }  
        },
        {
          name:"Personal Chat",
          action: () => {
            this.$refs.contextMenu.close();
            this.$store.dispatch("chat/startPersonalChat", {userId, username})
          }
        }
      ]
      
      this.$refs.contextMenu.open(e)
    },
    cmCleanup() {
      this.elementCfg = {};
      this.elementActions = [];
    }
  }
}
</script>

<style scoped>
.usersListItem {
  padding: 4px;
}
.usersListItem:hover {
  background-color: hsl(0, 0%, 95%);
}
</style>