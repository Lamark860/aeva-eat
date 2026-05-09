<template>
  <span class="sb-tape" :class="variant" :style="style"></span>
</template>

<script setup>
defineProps({
  variant: { type: String, default: '' },
  style:   { type: [Object, String], default: () => ({}) },
})
</script>
