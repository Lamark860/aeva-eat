<template>
  <span
    class="sb-author-tag"
    :class="[color, { 'has-photo': showPhoto }]"
    :style="positionStyle"
    :title="title"
  >
    <img
      v-if="showPhoto"
      :src="src"
      :alt="title || ''"
      class="ph"
      @error="onError"
    />
    <template v-else>{{ initial }}</template>
  </span>
</template>

<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  name:     { type: String, default: '' },
  color:    { type: String, default: 'terra' },
  title:    { type: String, default: '' },
  src:      { type: String, default: '' },                 // avatar_url, optional
  position: { type: Object, default: () => ({}) },
})

const broken = ref(false)
watch(() => props.src, () => { broken.value = false })

function onError() { broken.value = true }

const showPhoto = computed(() => !!props.src && !broken.value)

const initial = computed(() => (props.name || '?').slice(0, 1).toUpperCase())

const positionStyle = computed(() => {
  const p = props.position || {}
  const out = {}
  for (const key of ['top', 'right', 'bottom', 'left']) {
    if (p[key] !== undefined) out[key] = typeof p[key] === 'number' ? `${p[key]}px` : p[key]
  }
  return out
})
</script>

<style scoped>
/* The .sb-author-tag base styles live in scrapbook.scss; we only add the
   image fill here so consumers don't need to know whether avatar exists. */
.sb-author-tag.has-photo {
  /* Drop the colored fill so the image shows through cleanly. */
  background: var(--sb-paper-card);
  overflow: hidden;
}
.sb-author-tag .ph {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
