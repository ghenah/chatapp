<template>
  <div v-show="show" class="popup-window" :style="style" ref="context" tabindex="0" @blur="close">
    <slot>Pop-up window</slot>
  </div>
</template>

<script>
import {nextTick} from "vue";

export default {
  name: "BasePopupWindow",
  emits: ["closed"],
  props: ["display", "dimensions"],
  data() {
    return {
      // width: 400,
      // height: 650,
      show: false,
    }
  },
  computed: {
    left() {
      return window.innerWidth / 2 - this.dimensions.width / 2;
    },
    top() {
      return window.innerHeight / 2 - this.dimensions.height / 2;
    },
    style() {
      return {
        top: this.top + "px",
        left: this.left + "px",
        width: this.dimensions.width + "px",
        height: this.dimensions.height + "px",
      }
    }
  },
  methods: {
    close() {
      this.show = false;
      this.$emit("closed")
    },
    open() {
      nextTick(() => this.$el.focus());
      this.show = true;
    }
  }
}
</script>

<style scoped>
.popup-window {
  /* margin: 0 auto; */
  position: fixed;
  z-index: 990;
  outline: none;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.12), 0 1px 2px rgba(0, 0, 0, 0.24);
  cursor: pointer;
  background-color: hsl(230, 25%, 73%);
  overflow-y: scroll;
}
</style>