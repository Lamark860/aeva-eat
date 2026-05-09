<template>
  <router-link :to="`/places/${place.id}`" class="sb-artifact" :aria-label="place.name">
    <Polaroid
      :src="place.image_url || ''"
      :caption="caption"
      :gem="!!place.is_gem_place"
      :placeholder="placeholderClass"
    >
      <Tape :variant="tapeVariant" :style="tapeStyle" />
      <span v-if="place.is_gem_place" class="sb-gem-corner">
        <GemBadge :size="22" />
      </span>
      <AuthorTag
        v-if="firstAuthor"
        :name="firstAuthor.username"
        :color="firstAuthorColor"
        :title="firstAuthor.username"
        :src="firstAuthor.avatar_url || ''"
        :position="{ bottom: -6, left: -6 }"
      />
    </Polaroid>

    <div v-if="hasRatings" class="sb-ratings sb-t-l1">
      <Ticket
        compact
        :food="place.avg_food"
        :service="place.avg_service"
        :vibe="place.avg_vibe"
      />
    </div>

    <div v-if="metaLine" class="sb-meta">{{ metaLine }}</div>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'
import Polaroid from './Polaroid.vue'
import Tape from './Tape.vue'
import Ticket from './Ticket.vue'
import GemBadge from './GemBadge.vue'
import AuthorTag from './AuthorTag.vue'
import { authorColor, formatVisitCaption } from '../../composables/useFeed'

const props = defineProps({
  place: { type: Object, required: true },
})

const placeholderPalette = ['sb-photo-warm', 'sb-photo-olive', 'sb-photo-dusk', 'sb-photo-sage', 'sb-photo-peach', 'sb-photo-brick', 'sb-photo-cream', 'sb-photo-slate', 'sb-photo-indigo']
const placeholderClass = computed(() => placeholderPalette[(props.place.id ?? 0) % placeholderPalette.length])

const tapeVariants = ['', 'rose', 'mint', 'blue']
const tapeVariant = computed(() => tapeVariants[(props.place.id ?? 0) % tapeVariants.length])

// Tape position varies per item for a less-uniform feel.
const tapeStyle = computed(() => {
  const variants = [
    { top: '-10px', left: '32px', transform: 'rotate(-12deg)' },
    { top: '-8px',  left: '36px', transform: 'rotate(6deg)' },
    { top: '-8px',  right: '16px', transform: 'rotate(8deg)' },
    { top: '-9px',  right: '28px', transform: 'rotate(-6deg)' },
  ]
  return variants[(props.place.id ?? 0) % variants.length]
})

const caption = computed(() => {
  const d = props.place.created_at || props.place.updated_at
  if (!d) return props.place.name
  return formatVisitCaption(props.place.name, d)
})

const firstAuthor = computed(() => {
  const list = props.place.reviewers || []
  return list[0] || null
})
const firstAuthorColor = computed(() => authorColor(firstAuthor.value?.id ?? props.place.created_by))

const hasRatings = computed(() => {
  const p = props.place
  return [p.avg_food, p.avg_service, p.avg_vibe].some((v) => v !== null && v !== undefined)
})

const metaLine = computed(() => {
  const parts = []
  if (props.place.city) parts.push(props.place.city)
  if (props.place.cuisine_type) parts.push(props.place.cuisine_type)
  return parts.join(' · ')
})
</script>

<style scoped lang="scss">
.sb-artifact {
  display: block;
  text-decoration: none;
  color: inherit;
  position: relative;
  width: 100%;
}

.sb-gem-corner {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 2;
}

.sb-ratings {
  display: block;
  margin-top: 10px;
  margin-left: 4px;
}

.sb-meta {
  font-family: var(--sb-hand);
  font-size: 14px;
  color: var(--sb-ink-mute);
  margin-top: 4px;
  padding-left: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
