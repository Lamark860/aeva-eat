<template>
  <!-- Q-video: стопка видео-кружочков с inline-play. Используется и в фото-,
       и в ticket-only артефактах — DOM/логика общие, чтобы не дублировать
       код. Размер регулируется prop'ом `size` (normal/large). -->
  <div v-if="visibleVideos.length" class="art-kruzhoki" :class="`art-kruzhoki--${size}`">
    <span
      v-for="(url, i) in visibleVideos"
      :key="`kr-${i}-${url}`"
      class="art-kruzhok-layer"
      :style="kruzhokLayerStyle(i)"
      @click.stop.prevent="onKruzhokClick"
    >
      <video
        class="art-kruzhok-video"
        :src="url"
        preload="metadata"
        muted
        playsinline
        disablepictureinpicture
        aria-hidden="true"
        @loadedmetadata="forcePoster"
        @ended="onVideoEnded"
      ></video>
      <span class="art-kruzhok-play" aria-hidden="true"></span>
      <span class="art-kruzhok-hint" aria-hidden="true">тапни</span>
    </span>
    <span v-if="extraVideos > 0" class="art-kruzhok-extra">+{{ extraVideos }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  videos: { type: Array, default: () => [] },
  size:   { type: String, default: 'normal' }, // 'normal' | 'large'
})

const MAX_VIDEOS = 3
const visibleVideos = computed(() => props.videos.slice(0, MAX_VIDEOS))
const extraVideos = computed(() => Math.max(0, props.videos.length - MAX_VIDEOS))

function kruzhokLayerStyle(i) {
  const tilts = ['rotate(-5deg)', 'rotate(4deg)', 'rotate(-3deg)']
  return { top: `${i * 38}px`, transform: tilts[i % tilts.length], zIndex: i + 1 }
}

// preload=metadata в Safari/Chrome не всегда отрисовывает первый кадр — сикаем
// на 0.1s. Если за 500мс currentTime всё ещё 0 — показываем «тапни» fallback.
function forcePoster(ev) {
  const v = ev.target
  if (!v || v.currentTime > 0) return
  try { v.currentTime = 0.1 } catch (_) { /* fail-soft */ }
  setTimeout(() => {
    if (!v.isConnected || v.currentTime > 0) return
    v.closest('.art-kruzhok-layer')?.classList.add('posterless')
  }, 500)
}

function onKruzhokClick(ev) {
  const layer = ev.currentTarget
  const video = layer.querySelector('video')
  if (!video) return
  if (video.paused) {
    video.play().then(() => layer.classList.add('playing')).catch(() => {})
  } else {
    video.pause()
    layer.classList.remove('playing')
  }
}

function onVideoEnded(ev) {
  const v = ev.target
  if (!v) return
  v.currentTime = 0
  v.parentElement?.classList.remove('playing')
}
</script>

<style scoped lang="scss">
.art-kruzhoki {
  position: relative;
  flex: 0 0 auto;
  align-self: stretch;
}
.art-kruzhoki--normal { width: 64px; }
.art-kruzhoki--large  { width: 80px; }

.art-kruzhok-layer {
  position: absolute;
  border-radius: 50%;
  overflow: hidden;
  background: transparent;
  box-shadow:
    0 0 0 3px var(--sb-paper-card),
    0 2px 4px rgba(40, 30, 20, 0.18),
    0 6px 14px rgba(40, 30, 20, 0.16);
  cursor: pointer;
}
.art-kruzhoki--normal .art-kruzhok-layer { width: 56px; height: 56px; left: 4px; }
.art-kruzhoki--large  .art-kruzhok-layer { width: 70px; height: 70px; left: 5px; }

.art-kruzhok-video {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  border-radius: 50%;
}

.art-kruzhok-play {
  position: absolute;
  inset: 0;
  pointer-events: none;
  transition: opacity 180ms ease;
  z-index: 2;
}
.art-kruzhok-play::before,
.art-kruzhok-play::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
}
.art-kruzhok-play::before {
  width: 28px;
  height: 28px;
  margin: -14px 0 0 -14px;
  border-radius: 50%;
  background: rgba(20, 12, 6, 0.55);
  backdrop-filter: blur(2px);
  -webkit-backdrop-filter: blur(2px);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.35);
}
.art-kruzhok-play::after {
  width: 0;
  height: 0;
  border-style: solid;
  border-width: 5px 0 5px 8px;
  border-color: transparent transparent transparent #fff;
  margin: -5px 0 0 -3px;
  filter: drop-shadow(0 1px 0.5px rgba(0, 0, 0, 0.25));
}
.art-kruzhoki--large .art-kruzhok-play::before {
  width: 32px;
  height: 32px;
  margin: -16px 0 0 -16px;
}
.art-kruzhoki--large .art-kruzhok-play::after {
  border-width: 6px 0 6px 9px;
  margin: -6px 0 0 -3.5px;
}
.art-kruzhok-layer.playing .art-kruzhok-play { opacity: 0; }
.art-kruzhok-layer.playing .art-kruzhok-hint { opacity: 0; }

.art-kruzhok-hint {
  position: absolute;
  inset: 0;
  display: none;
  align-items: center;
  justify-content: center;
  font-family: var(--sb-hand);
  font-size: 14px;
  color: var(--sb-paper-card);
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.6);
  pointer-events: none;
  z-index: 3;
}
.art-kruzhok-layer.posterless {
  background: oklch(0.35 0.04 60);
}
.art-kruzhok-layer.posterless .art-kruzhok-hint { display: flex; }
.art-kruzhok-layer.posterless .art-kruzhok-play { opacity: 0; }

.art-kruzhok-extra {
  position: absolute;
  bottom: -14px;
  right: -2px;
  font-family: var(--sb-hand);
  font-size: 14px;
  color: var(--sb-ink-mute);
  transform: rotate(-3deg);
  pointer-events: none;
}
</style>
