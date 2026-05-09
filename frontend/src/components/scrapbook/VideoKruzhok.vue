<template>
  <div class="sb-kruzhok" :style="{ width: size + 'px', height: size + 'px' }">
    <div class="play"></div>
    <div v-if="duration" class="duration">{{ duration }}</div>
    <slot />
  </div>
</template>

<script setup>
defineProps({
  duration: { type: String, default: '' },
  size:     { type: Number, default: 86 },
})
</script>
