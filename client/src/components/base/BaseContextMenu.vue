<template>
  <div v-show="show" class="context-menu" :style="style" ref="context" tabindex="0" @blur="close">
    <p class="menuName" v-if="cfg.hasName">
        <slot name="menuName" :cmName="cfg.name"></slot>
    </p>
    <ul id=actions>
      <li v-for="(action, i) in actions" :key="i" class="itemAction">
        <slot name="action" :actionName="action.name" :cb="action.action"></slot>
      </li>
    </ul>
  </div>
</template>


<script>
import {nextTick} from "vue";

export default {
  name: "BaseContextMenu",
  emits: ["closed"],
  props: ["display", "cfg","actions"],
  data() {
    return {
      left: 0,
      top: 0,
      show: false,
    }
  },
  computed: {
    style() {
      return {
        top: this.top + "px",
        left: this.left + "px",
      }
    }
  },
  methods: {
    close() {
      this.show = false;
      this.left = 0;
      this.top = 0;
      this.$emit("closed")
    },
    open(e) {
      this.left = e.pageX || e.clientX;
      this.top = e.pageY || e.clientY;

      nextTick(() => this.$el.focus());
      this.show = true;
    }
  }
}
</script>

<style scoped>
.context-menu {
  position: fixed;
  z-index: 999;
  outline: none;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.12), 0 1px 2px rgba(0, 0, 0, 0.24);
  cursor: pointer;
}
ul, p {
  margin: 0;
  padding: 0;
}
.menuName {
  padding: 4px;
  color: darkblue;
  font-weight: bold;
  background-color: hsl(230, 30%, 70%);
}
.itemAction {
  padding: 4px;
  color: hsl(230, 5%, 100%);
  background-color: hsl(230, 25%, 73%);
}
.itemAction:hover {
  background-color: hsl(230, 30%, 70%);
}

</style>