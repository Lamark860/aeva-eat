<template>
  <!--
    NEXT.md §A2 — разделитель-строка с кружочком видео между артефактами.
    Слева/справа — тонкие чернильные линии, в центре — кружочек, справа дата
    рукописным. Клик ведёт на страницу места с anchor'ом на отзыв.
  -->
  <router-link :to="targetPath" class="sb-kruzhok-divider" :aria-label="`видео из ${placeName}`">
    <span class="ink-line left" aria-hidden="true"></span>
    <span class="kruzhok-wrap">
      <span class="kruzhok-frame">
        <span class="play" aria-hidden="true">▶</span>
      </span>
    </span>
    <span class="ink-line right" aria-hidden="true"></span>
    <span v-if="dateLabel" class="date">{{ dateLabel }}</span>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  placeId:   { type: [Number, String], required: true },
  placeName: { type: String, default: '' },
  date:      { type: String, default: '' },
})

const targetPath = computed(() => `/places/${props.placeId}#detail-review-form`)

const dateLabel = computed(() => {
  if (!props.date) return ''
  const d = new Date(props.date)
  if (Number.isNaN(d.getTime())) return ''
  return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' })
})
</script>

<style scoped lang="scss">
.sb-kruzhok-divider {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 4px;
  text-decoration: none;
  color: inherit;
}

.ink-line {
  flex: 1 1 auto;
  height: 1.5px;
  background: rgba(40, 30, 20, 0.22);
  border-radius: 1px;

  &.left  { transform: rotate(-0.4deg); }
  &.right { transform: rotate(0.4deg); }
}

.kruzhok-wrap {
  flex: 0 0 auto;
}

.kruzhok-frame {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: oklch(0.32 0.04 50);
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  /* «Скрапбук-обводка» как у VideoKruzhok — paper ring + soft shadow. */
  box-shadow:
    0 0 0 3px var(--sb-paper-card),
    0 2px 4px rgba(40, 30, 20, 0.18),
    0 6px 14px rgba(40, 30, 20, 0.16);
  transition: transform 200ms ease;
}
.sb-kruzhok-divider:hover .kruzhok-frame {
  transform: scale(1.04);
}

.play {
  color: var(--sb-paper-card);
  font-size: 18px;
  line-height: 1;
  transform: translateX(1px);
}

.date {
  flex: 0 0 auto;
  font-family: var(--sb-hand);
  font-size: 14px;
  color: var(--sb-ink-mute);
  white-space: nowrap;
}
</style>
