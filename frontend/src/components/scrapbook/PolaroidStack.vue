<template>
  <!--
    Стопка из 1–3 видимых полароидов по DESIGN-DECISIONS §L3:
    – 1 фото:  одиночный полароид (как Polaroid.vue)
    – 2 фото:  два внахлёст со встречным наклоном
    – 3+ фото: стопка из трёх с каракулевой пометкой "+ ещё N"
    Слот пробрасывается на верхний полароид (там же сидят Tape / GemBadge / AuthorTag).
  -->
  <div class="sb-polaroid-stack" :class="{ single: visible.length === 1 }" :style="rootStyle">
    <Polaroid
      v-for="(p, i) in visible"
      :key="p.key"
      :src="p.url"
      :placeholder="placeholder"
      class="layer"
      :class="[`layer-${i}`, layerTilt(i)]"
      :style="layerStyle(i)"
      :gem="gem && i === topIndex"
    >
      <slot v-if="i === topIndex" />
    </Polaroid>

    <span v-if="extraCount > 0" class="extra">+ ещё {{ extraCount }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import Polaroid from './Polaroid.vue'

const props = defineProps({
  // Источник — либо массив объектов { id, url }, либо просто массив url-строк.
  // Один url-string или null/undefined → пустая стопка (рисуем placeholder).
  photos:      { type: Array, default: () => [] },
  placeholder: { type: String, default: 'sb-photo-cream' },
  gem:         { type: Boolean, default: false },
  // Опциональный fallback: если photos пустой, всё равно показать что-то.
  // Используется в ленте, где есть place.image_url как cover.
  fallbackSrc: { type: String, default: '' },
})

// Нормализуем вход: { id, url } или строка → { url, key }
const normalized = computed(() => {
  const list = (props.photos || []).filter(Boolean).map((p, i) => {
    if (typeof p === 'string') return { url: p, key: `s-${i}-${p}` }
    return { url: p.url, key: p.id ?? `i-${i}-${p.url}` }
  })
  if (list.length === 0 && props.fallbackSrc) {
    return [{ url: props.fallbackSrc, key: 'fallback' }]
  }
  return list
})

// Показываем максимум три верхних — четвёртый и далее свернуты в caption.
const visible = computed(() => normalized.value.slice(0, 3))
const extraCount = computed(() => Math.max(0, normalized.value.length - 3))
const topIndex = computed(() => Math.max(0, visible.value.length - 1))

// Слой 0 — самый нижний; верхний (последний) — это «фокусное» фото.
// Наклоны: для 2 — встречные −2°/+3°, для 3 — мягкая стопка.
function layerTilt() { return '' /* tilt задаём через inline transform */ }

function layerStyle(i) {
  const total = visible.value.length
  if (total === 1) return {}
  if (total === 2) {
    // два внахлёст: первый чуть смещён вниз-влево, второй — наверх-вправо
    if (i === 0) {
      return {
        transform: 'translate(-6%, 4%) rotate(-2deg)',
        zIndex: 1,
      }
    }
    return {
      transform: 'translate(6%, -2%) rotate(3deg)',
      zIndex: 2,
    }
  }
  // total === 3
  if (i === 0) return { transform: 'translate(-8%, 6%) rotate(-4deg)', zIndex: 1 }
  if (i === 1) return { transform: 'translate(0, 0) rotate(2deg)',     zIndex: 2 }
  return         { transform: 'translate(7%, -3%) rotate(-1deg)',      zIndex: 3 }
}

const rootStyle = computed(() => {
  // Стопка отъедает чуть больше места по бокам и снизу-сверху, чтобы
  // повёрнутые слои не клипались родителем.
  if (visible.value.length === 1) return {}
  return {
    paddingInline: '8%',
    paddingBlock: '4% 6%',
  }
})
</script>

<style scoped lang="scss">
.sb-polaroid-stack {
  position: relative;
  width: 100%;

  &.single {
    /* В single-режиме рендер идентичен одиночному полароиду */
  }
}

.layer {
  /* Когда стопка из 2+ — слои абсолютно позиционируются друг над другом */
  transition: transform 220ms ease;
}

.sb-polaroid-stack:not(.single) .layer {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  /* Каждый слой получает margin от родительского padding (см. rootStyle).
     inset: 0 ставит его в логическую плоскость, а transform даёт смещение. */
}

.sb-polaroid-stack:not(.single) {
  /* Чтобы родителю было что измерять — фиксируем aspect-ratio квадрата
     (полароид fluid делает то же самое). */
  aspect-ratio: 1 / 1;
}

.extra {
  position: absolute;
  right: 4px;
  bottom: -10px;
  font-family: var(--sb-hand);
  font-size: 16px;
  color: var(--sb-ink-mute);
  /* «каракулями, не круглый бейдж» — лёгкий наклон, без подложки */
  transform: rotate(-3deg);
  z-index: 5;
  pointer-events: none;
  white-space: nowrap;
}
</style>
