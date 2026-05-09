<template>
  <span class="sb-ticket" :class="{ compact }">
    <span v-for="s in stubs" :key="s.lbl" class="stub">
      <span class="lbl">{{ s.lbl }}</span>
      <span class="val">{{ s.val }}</span>
    </span>
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  food:    { type: [Number, String, null], default: null },
  service: { type: [Number, String, null], default: null },
  vibe:    { type: [Number, String, null], default: null },
  labels:  { type: Array,   default: () => ['еда', 'сервис', 'вайб'] },
  compact: { type: Boolean, default: false },
})

function fmt(v) {
  if (v === null || v === undefined || v === '') return '–'
  const n = Number(v)
  if (Number.isNaN(n)) return String(v)
  return Number.isInteger(n) ? String(n) : n.toFixed(1)
}

const stubs = computed(() => [
  { lbl: props.labels[0] || 'еда',    val: fmt(props.food) },
  { lbl: props.labels[1] || 'сервис', val: fmt(props.service) },
  { lbl: props.labels[2] || 'вайб',   val: fmt(props.vibe) },
])
</script>
