<template>
  <div>
    <create-user @create-user-submitted="createNewUser"></create-user>
    <get-user @get-user-by-id="getUserById" :user-id="userId" :username="username" :email="email"></get-user>
  </div>
</template>

<script>
  import CreateUser from "./components/CreateUser"
  import GetUser from "./components/GetUser"
  import axios from "axios"
export default {
  name: "Golang-docker",
  data() {
    return {

    userId: "",
    username: "",
    email: ""
    }
  },
  components: {
    CreateUser,
    GetUser,
  },
  methods: {
    createNewUser(id, username, email) {
      const res = axios.post("http://localhost:8001/create", {
        id: parseInt(id),
        username,
        email
      }).then(
        console.log(res)
      )
    },
    getUserById(id) {
      axios.get("http://localhost:8001/get/"+parseInt(id))
      .then(res => {
        console.log(res)
        console.log(res.data.id)
        console.log(res.data.username)
        console.log(res.data.email)
        this.userId = res.data.id
        this.username = res.data.username
        this.email = res.data.email
      }
      )
    }
  }
}
</script>