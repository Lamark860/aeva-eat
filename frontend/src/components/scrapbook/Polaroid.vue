<template>
  <div
    class="sb-polaroid"
    :class="[gem ? 'gem-frame' : '', tilt ? `sb-${tilt}` : '', fluid ? 'fluid' : '']"
  >
    <img
      v-if="src && !broken"
      class="photo"
      :src="src"
      :alt="caption || ''"
      :style="photoStyle"
      loading="lazy"
      @error="onError"
    />
    <div
      v-else
      class="photo"
      :class="placeholderClass"
      :style="photoStyle"
    ></div>
    <div v-if="caption" class="caption">{{ caption }}</div>
    <slot />
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  src:         { type: String, default: '' },
  caption:     { type: String, default: '' },
  width:       { type: Number, default: 0 },                            // 0 → fluid (100% of column)
  height:      { type: Number, default: 0 },
  tilt:        { type: String, default: '' },                           // e.g. 't-l3', 't-r2'
  gem:         { type: Boolean, default: false },
  placeholder: { type: String, default: 'sb-photo-cream' },             // SCSS class for missing image
})

const broken = ref(false)
watch(() => props.src, () => { broken.value = false })

function onError() { broken.value = true }

const placeholderClass = computed(() => props.placeholder)
const fluid = computed(() => !props.width)

const photoStyle = computed(() => {
  if (fluid.value) return { width: '100%', aspectRatio: '1 / 1' }
  return {
    width: `${props.width}px`,
    height: `${props.height || props.width}px`,
  }
})
</script>
