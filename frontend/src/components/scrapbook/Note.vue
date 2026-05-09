<template>
  <div class="sb-note" :class="{ lined }">
    <slot />
  </div>
</template>

<script setup>
defineProps({
  lined: { type: Boolean, default: false },
})
</script>
