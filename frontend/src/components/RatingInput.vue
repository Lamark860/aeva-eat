<template>
  <div class="sb-rating" :data-tone="tone">
    <div class="head">
      <span class="lbl">{{ label }}</span>
      <span class="val">{{ display }}</span>
    </div>
    <div class="track-wrap">
      <div class="track-line"></div>
      <div class="track-fill" :style="{ width: fillPercent + '%' }"></div>
      <input
        type="range"
        min="0"
        max="10"
        step="0.1"
        :value="modelValue"
        class="track-input"
        @input="$emit('update:modelValue', parseFloat($event.target.value))"
      />
    </div>
    <div class="ticks">
      <span>0</span>
      <span>5</span>
      <span>10</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: { type: Number, default: 0 },
  label:      { type: String, default: '' },
  /* tone determines the fill color: terra (default), ochre, moss, plum */
  tone:       { type: String, default: 'terra' },
})

defineEmits(['update:modelValue'])

const fillPercent = computed(() => (props.modelValue ? (props.modelValue / 10) * 100 : 0))

const display = computed(() => {
  const v = props.modelValue
  if (v === null || v === undefined || v === 0) return '–'
  return Number(v).toFixed(1)
})
</script>

<style scoped lang="scss">
.sb-rating {
  margin-bottom: 6px;
}

.head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-bottom: 2px;
}
.lbl {
  font-family: var(--sb-serif);
  font-style: italic;
  font-weight: 500;
  font-size: 14px;
  color: var(--sb-ink);
}
.val {
  font-family: var(--sb-hand);
  font-size: 22px;
  line-height: 1;
  color: var(--sb-ink);
}

.track-wrap {
  position: relative;
  height: 32px;
  display: flex;
  align-items: center;
}

/* Thin ink line as the track */
.track-line {
  position: absolute;
  left: 0;
  right: 0;
  height: 1.5px;
  background: rgba(40, 30, 20, 0.35);
  border-radius: 1px;
  pointer-events: none;
}
.track-fill {
  position: absolute;
  left: 0;
  height: 2.5px;
  background: var(--fill-color);
  border-radius: 2px;
  pointer-events: none;
  transition: width 0.15s ease;
}

.track-input {
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

/* Thumb — pushpin head */
.track-input::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: radial-gradient(circle at 35% 30%, oklch(0.7 0.16 25), oklch(0.42 0.18 25));
  box-shadow:
    inset 0 -1px 1px rgba(0, 0, 0, 0.3),
    inset 0 1px 1px rgba(255, 255, 255, 0.3),
    0 1px 2px rgba(40, 15, 5, 0.3),
    0 0 0 3px var(--sb-paper-card);
  cursor: pointer;
  transition: transform 0.15s;
}
.track-input::-webkit-slider-thumb:hover {
  transform: scale(1.12);
}
.track-input::-moz-range-thumb {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: radial-gradient(circle at 35% 30%, oklch(0.7 0.16 25), oklch(0.42 0.18 25));
  box-shadow:
    inset 0 -1px 1px rgba(0, 0, 0, 0.3),
    inset 0 1px 1px rgba(255, 255, 255, 0.3),
    0 1px 2px rgba(40, 15, 5, 0.3),
    0 0 0 3px var(--sb-paper-card);
  border: none;
  cursor: pointer;
}

.ticks {
  display: flex;
  justify-content: space-between;
  font-family: var(--sb-hand);
  font-size: 13px;
  color: var(--sb-ink-mute);
  line-height: 1;
  margin-top: -2px;
}

/* Tone-driven fill colours via CSS var so we don't repeat selectors */
.sb-rating { --fill-color: var(--sb-terracotta); }
.sb-rating[data-tone='ochre'] { --fill-color: var(--sb-ochre); }
.sb-rating[data-tone='moss']  { --fill-color: var(--sb-moss); }
.sb-rating[data-tone='plum']  { --fill-color: var(--sb-plum); }
</style>
