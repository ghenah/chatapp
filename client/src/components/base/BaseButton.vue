<template>
<div role="button" v-if="btnType === 'default'" :class="btnClass">
  <slot name="default"></slot>
</div>
</template>

<script>
export default {
  props: {
    btnType: String,
    btnClass: String,
  },
}
</script>

<style scoped>
.default {
  display: inline-block;
  background-color: #ffffff;
  border: 1px solid #181818;
  padding: 2px;
}
.positive {
  display: inline-block;
  background-color: #39ac73;
  padding: 2px;
}
.negative {
  display: inline-block;
  background-color: #bf4040;
  padding: 2px;
}
</style>