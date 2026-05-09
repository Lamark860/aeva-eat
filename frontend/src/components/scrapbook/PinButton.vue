<template>
  <button class="sb-pin-btn" type="button" @click="emit('click', $event)">
    <span class="plus">+</span>
    <span><slot>добавить</slot></span>
  </button>
</template>

<script setup>
const emit = defineEmits(['click'])
</script>
