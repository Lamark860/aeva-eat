<template>
  <div class="sb-collapsed-strip">
    <div class="head">
      <div class="dates">{{ dates }}</div>
      <div class="rule"></div>
      <div class="meta">
        {{ count }} {{ countLabel }}
        <template v-if="gemCount > 0">
          · {{ gemCount }} <span style="color: var(--sb-terracotta)">♦</span>
        </template>
      </div>
    </div>
    <div class="body">
      <div class="stack">
        <div
          v-for="(p, i) in summary"
          :key="p.id ?? i"
          class="sb-mini"
          :style="miniStyle(i)"
        >
          <div class="sb-polaroid" style="padding: 5px 5px 18px">
            <img
              v-if="p.src"
              class="photo"
              :src="p.src"
              alt=""
              loading="lazy"
              style="width: 60px; height: 60px"
            />
            <div
              v-else
              class="photo"
              :class="p.placeholder || pickPlaceholder(i)"
              style="width: 60px; height: 60px"
            ></div>
            <div class="caption" style="font-size: 11px; bottom: 3px">{{ p.cap }}</div>
            <div v-if="p.gem" style="position: absolute; top: 4px; right: 4px">
              <GemBadge :size="16" />
            </div>
          </div>
        </div>
      </div>
      <button class="expand-btn" type="button" @click="emit('expand')">
        раскрыть&nbsp;↓
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import GemBadge from './GemBadge.vue'

const props = defineProps({
  dates:    { type: String, required: true },
  summary:  { type: Array,  required: true },          // [{ id, src, cap, gem, placeholder }]
  count:    { type: Number, required: true },
  gemCount: { type: Number, default: 0 },
})

const emit = defineEmits(['expand'])

const placeholders = ['sb-photo-warm', 'sb-photo-olive', 'sb-photo-dusk', 'sb-photo-sage', 'sb-photo-peach', 'sb-photo-indigo']
function pickPlaceholder(i) { return placeholders[i % placeholders.length] }

const countLabel = computed(() => {
  const n = props.count % 100
  if (n >= 11 && n <= 14) return 'мест'
  const last = n % 10
  if (last === 1) return 'место'
  if (last >= 2 && last <= 4) return 'места'
  return 'мест'
})

function miniStyle(i) {
  return {
    position: 'absolute',
    left: `${i * 46}px`,
    top: i % 2 === 0 ? '4px' : '12px',
    transform: `rotate(${(i % 2 === 0 ? -2 : 2) - i * 0.3}deg)`,
  }
}
</script>
