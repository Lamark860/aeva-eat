<template>
  <div class="rating-slider">
    <div class="d-flex align-items-center justify-content-between mb-1">
      <span class="rating-label">{{ label }}</span>
      <span class="rating-value" :style="{ color: fillColor }">{{ modelValue !== null && modelValue !== undefined ? Number(modelValue).toFixed(1) : '–' }}</span>
    </div>
    <div class="slider-track-wrap" @click="onTrackClick" ref="trackEl">
      <div class="slider-track">
        <div class="slider-fill" :style="{ width: fillPercent + '%', background: fillColor }"></div>
      </div>
      <input
        type="range"
        min="0"
        max="10"
        step="0.1"
        :value="modelValue"
        class="slider-input"
        @input="$emit('update:modelValue', parseFloat($event.target.value))"
      />
    </div>
    <div class="d-flex justify-content-between text-muted" style="font-size: 0.65rem;">
      <span>0</span>
      <span>5</span>
      <span>10</span>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  modelValue: { type: Number, default: 0 },
  label: { type: String, default: '' },
  color: { type: String, default: '#f59e0b' }
})

defineEmits(['update:modelValue'])

const trackEl = ref(null)

const fillPercent = computed(() => props.modelValue ? (props.modelValue / 10) * 100 : 0)
const fillColor = computed(() => props.color)

function onTrackClick(e) {
  // handled by range input
}
</script>

<style scoped>
.rating-slider {
  margin-bottom: 0.25rem;
}

.rating-label {
  font-size: 0.85rem;
  font-weight: 500;
  color: #666;
}

.rating-value {
  font-size: 1.1rem;
  font-weight: 700;
}

.slider-track-wrap {
  position: relative;
  height: 32px;
  display: flex;
  align-items: center;
}

.slider-track {
  position: absolute;
  left: 0; right: 0;
  height: 6px;
  background: #e8e6e3;
  border-radius: 3px;
  overflow: hidden;
  pointer-events: none;
}

.slider-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.15s ease;
}

.slider-input {
  -webkit-appearance: none;
  appearance: none;
  width: 100%;
  height: 32px;
  background: transparent;
  position: relative;
  z-index: 2;
  cursor: pointer;
  margin: 0;
}

.slider-input::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0,0,0,0.2), 0 0 0 2px v-bind(fillColor);
  cursor: pointer;
  transition: transform 0.15s;
}
.slider-input::-webkit-slider-thumb:hover {
  transform: scale(1.15);
}

.slider-input::-moz-range-thumb {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0,0,0,0.2), 0 0 0 2px v-bind(fillColor);
  border: none;
  cursor: pointer;
}
</style>
