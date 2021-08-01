<template>
  <div class="profile-pic-container">
    <img :src="'/images/profile/' + profilePicture" alt="Profile picture" class="profile-pic">
    <div class="dropzone">
      <input type="file" id="dropzone-file" class="dropzone-file" @change="handleUpload" ref="dropzoneFile" accept="image/*">
      <div class="dropzone-wrapper"  @drop.prevent="handleUpload($event)" @dragenter.prevent="" @dragover.prevent="">
        <label for="dropzone-file" class="dropzone-label">
          <div class="dropzone-label-text">Drag here to update</div>
          <div role="button" class="dropzone-select-button">or select manually</div>
        </label>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "ChangePicture",
  data() {
    return {
      file: null
    }
  },
  computed: {
    profilePicture() {
      return this.$store.getters["user/profilePicture"];
    }
  },
  methods: {
    handleUpload(e) {
      const inputValue = e.target.files || e.dataTransfer.files || this.$refs.dropzoneFile.files;
      const profilePic = inputValue[0];

      const formData = new FormData();
      formData.append("mimeType", profilePic.type);
      formData.append("profilePic", profilePic);

      this.$store.dispatch("user/changeProfilePicture", {formData});
    }
  }
}
</script>

<style scoped>
.profile-pic-container {
  position: relative;
  height: 200px;
  width: 200px;
}
.profile-pic {
  z-index: 19;
  position: absolute;
  height: 200px;
  width: 200px;
}
.dropzone {
  z-index: 20;
  background-color: hsla(0, 0%, 0%, 30%);
  box-sizing: border-box;
  position: absolute;
  width: 200px;
  height: 200px;
  border: 2px dashed hsla(180, 100%, 50%, 50%);
  text-align: center;
}
#dropzone-file {
  display: none;
}
.dropzone-wrapper {
  margin: 0px;
  height: 100%;
  width: 100%;
}
.dropzone-label-text {
  margin-top: 10px;
  color: cyan;
  font-family: sans-serif;
  font-weight: bold;
  text-shadow: 1px 1px 2px black, 0 0 1em black, 0 0 0.2em black;
}
.dropzone-select-button {
  margin: auto;
  margin-top: 120px;
  display: inline-block;
  background-color: cyan;
  padding: 4px;
  border-radius: 10px;
}
</style>