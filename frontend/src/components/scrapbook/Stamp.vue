<template>
  <span class="sb-stamp" :class="kind">
    <slot />
  </span>
</template>

<script setup>
defineProps({
  kind: { type: String, default: '' },
})
</script>
