<template>
  <div class="rating-input d-flex align-items-center gap-1">
    <span class="text-muted me-1" style="min-width: 70px; font-size: 0.85rem;">{{ label }}</span>
    <button
      v-for="n in 10"
      :key="n"
      type="button"
      class="btn btn-sm rating-btn"
      :class="n <= modelValue ? activeClass : 'btn-outline-secondary'"
      @click="$emit('update:modelValue', n)"
      style="width: 30px; height: 30px; padding: 0; font-size: 0.75rem;"
    >
      {{ n }}
    </button>
    <span class="ms-2 fw-bold" :class="textClass">{{ modelValue || '–' }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: { type: Number, default: 0 },
  label: { type: String, default: '' },
  color: { type: String, default: 'warning' }
})

defineEmits(['update:modelValue'])

const activeClass = computed(() => `btn-${props.color}`)
const textClass = computed(() => `text-${props.color}`)
</script>

<style scoped>
.rating-btn {
  transition: all 0.15s;
}
.rating-btn:hover {
  transform: scale(1.15);
}
</style>
