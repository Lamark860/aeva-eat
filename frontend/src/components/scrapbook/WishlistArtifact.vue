<template>
  <router-link
    :to="`/places/${entry.place.id}`"
    class="wish-art"
    :class="{ struck: entry.is_struck }"
  >
    <Stamp v-if="!entry.is_struck" class="wish-stamp">план</Stamp>
    <Stamp v-else class="wish-stamp" kind="moss">сходили ✓</Stamp>

    <span v-if="!entry.is_struck" class="wish-pin" aria-hidden="true"></span>

    <div class="wish-title-wrap">
      <span class="wish-title">{{ entry.place.name }}</span>
      <!-- Q1: SVG-волнистый штрих, НЕ CSS line-through. Маркер по бумаге.
           viewBox растягивается шириной title через preserveAspectRatio=none. -->
      <svg
        v-if="entry.is_struck"
        class="wish-strike"
        viewBox="0 0 100 6"
        preserveAspectRatio="none"
        aria-hidden="true"
      >
        <path
          d="M 0 3 Q 5 0.5 10 3 T 20 3 T 30 3 T 40 3 T 50 3 T 60 3 T 70 3 T 80 3 T 90 3 T 100 3"
          fill="none"
          stroke="currentColor"
          stroke-width="1.6"
          stroke-linecap="round"
        />
      </svg>
    </div>

    <div v-if="entry.place.city" class="wish-city">{{ entry.place.city }}</div>

    <div class="wish-by">{{ entry.username }}</div>

    <!-- Q1: мини-полароид визита «внахлёст» — для зачёркнутых. -->
    <div
      v-if="entry.is_struck && entry.place.image_url"
      class="wish-mini-polaroid"
    >
      <img :src="entry.place.image_url" alt="" />
    </div>
  </router-link>
</template>

<script setup>
import Stamp from './Stamp.vue'

defineProps({
  entry: { type: Object, required: true },
})
</script>

<style scoped lang="scss">
.wish-art {
  display: block;
  position: relative;
  background: var(--sb-paper);
  padding: 14px 14px 16px;
  border-radius: 1px;
  text-decoration: none;
  color: var(--sb-ink);
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.08),
    0 4px 10px rgba(40, 30, 20, 0.10);
  // Лёгкая «бумажная» текстура — поддерживает контраст с другими карточками.
  min-height: 90px;

  &.struck {
    background: var(--sb-paper-deep);
    color: var(--sb-ink-soft);
  }
}

.wish-stamp {
  position: absolute;
  top: 8px;
  right: 8px;
  transform: rotate(4deg);
}

// Канцелярская кнопка в верхнем-левом углу — только для активных планов.
.wish-pin {
  position: absolute;
  top: -4px;
  left: 14px;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: radial-gradient(circle at 35% 30%,
    oklch(0.7 0.16 25),
    oklch(0.42 0.18 25));
  box-shadow:
    inset 0 -1px 1px rgba(0, 0, 0, 0.3),
    inset 0 1px 1px rgba(255, 255, 255, 0.3),
    0 1px 2px rgba(40, 15, 5, 0.4);
}

.wish-title-wrap {
  position: relative;
  display: inline-block;
  margin-top: 8px;
}
.wish-title {
  font-family: var(--sb-serif);
  font-style: italic;
  font-weight: 500;
  font-size: 18px;
  line-height: 1.15;
  color: inherit;
  word-break: break-word;
}
.wish-strike {
  position: absolute;
  left: -2px;
  right: -2px;
  top: 50%;
  width: calc(100% + 4px);
  height: 6px;
  transform: translateY(-50%);
  color: var(--sb-terracotta);
  pointer-events: none;
}

.wish-city {
  font-family: var(--sb-hand);
  font-size: 14px;
  color: var(--sb-ink-mute);
  margin-top: 4px;
}

.wish-by {
  font-family: var(--sb-hand);
  font-size: 13px;
  color: var(--sb-ink-mute);
  margin-top: 6px;
  font-style: italic;
}

.wish-mini-polaroid {
  position: absolute;
  bottom: -10px;
  right: -8px;
  width: 56px;
  height: 56px;
  background: var(--sb-paper-card);
  padding: 4px 4px 14px;
  transform: rotate(4deg);
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.08),
    0 4px 10px rgba(40, 30, 20, 0.14);
  border-radius: 1px;
  z-index: 2;

  img {
    display: block;
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}
</style>
