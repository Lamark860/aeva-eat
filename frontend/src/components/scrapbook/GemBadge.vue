<template>
  <span class="sb-gem-badge" :style="{ width: size + 'px', height: size + 'px' }">
    <svg viewBox="0 0 26 26">
      <defs>
        <linearGradient :id="gradId" x1="0" y1="0" x2="1" y2="1">
          <stop offset="0%" stop-color="oklch(0.95 0.04 85)" />
          <stop offset="100%" stop-color="oklch(0.85 0.07 60)" />
        </linearGradient>
      </defs>
      <path d="M13 2 L23 10 L13 24 L3 10 Z" :fill="`url(#${gradId})`" stroke="oklch(0.55 0.14 30)" stroke-width="1.6" stroke-linejoin="round" />
      <path d="M13 2 L8 10 L13 24 M13 2 L18 10 L13 24 M3 10 L23 10" stroke="oklch(0.55 0.14 30)" stroke-width="1" fill="none" opacity="0.6" />
      <g class="gleam" style="mix-blend-mode: screen">
        <path d="M9 5 L11 5 L7 11 L5 11 Z" fill="white" opacity="0.9" />
        <circle cx="11" cy="6" r="1.1" fill="white" opacity="0.95" />
      </g>
    </svg>
  </span>
</template>

<script setup>
import { computed } from 'vue'

defineProps({
  size: { type: Number, default: 26 },
})

// Unique gradient id per instance — required for multiple gems on one page.
const gradId = computed(() => `sb-gem-${Math.random().toString(36).slice(2, 8)}`)
</script>
