<template>
  <router-link :to="`/places/${place.id}`" class="sb-result">
    <div class="thumb">
      <Polaroid
        :src="place.image_url || ''"
        :width="76"
        :height="76"
        :tilt="tilt"
        :gem="!!place.is_gem_place"
        :placeholder="placeholderClass"
      >
        <Tape :variant="tapeVariant" :style="tapeStyle" />
        <span v-if="place.is_gem_place" class="gem-corner">
          <GemBadge :size="16" />
        </span>
      </Polaroid>
    </div>

    <div class="body">
      <div class="title-row">
        <h3 class="title">{{ place.name }}</h3>
        <Stamp v-if="place.is_gem_place" kind="gem" class="gem-stamp">♦</Stamp>
      </div>

      <div v-if="metaLine" class="meta">{{ metaLine }}</div>

      <div class="bottom">
        <Ticket
          v-if="hasRatings"
          compact
          class="mini-ticket"
          :food="place.avg_food"
          :service="place.avg_service"
          :vibe="place.avg_vibe"
        />
        <div v-if="reviewers.length" class="people">
          <span
            v-for="r in reviewers"
            :key="r.id"
            class="r-tag sb-author-tag"
            :class="[authorColor(r.id), { 'has-photo': r.avatar_url }]"
            :title="r.username"
          >
            <img v-if="r.avatar_url" :src="r.avatar_url" alt="" class="r-ph" />
            <template v-else>{{ (r.username || '?').slice(0, 1).toUpperCase() }}</template>
          </span>
        </div>
      </div>
    </div>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'
import Polaroid from './Polaroid.vue'
import Tape from './Tape.vue'
import Ticket from './Ticket.vue'
import Stamp from './Stamp.vue'
import GemBadge from './GemBadge.vue'
import { authorColor } from '../../composables/useFeed'

const props = defineProps({
  place: { type: Object, required: true },
})

const tilts = ['t-l3', 't-r2', 't-l2', 't-r3']
const tilt = computed(() => tilts[(props.place.id ?? 0) % tilts.length])

const tapeVariants = ['', 'rose', 'mint', 'blue']
const tapeVariant = computed(() => tapeVariants[(props.place.id ?? 0) % tapeVariants.length])

const tapeStyle = computed(() => ({
  top: '-8px',
  left: '20px',
  transform: 'rotate(-10deg)',
  width: '40px',
}))

const placeholderPalette = ['sb-photo-warm', 'sb-photo-olive', 'sb-photo-dusk', 'sb-photo-sage', 'sb-photo-peach', 'sb-photo-brick', 'sb-photo-cream', 'sb-photo-slate', 'sb-photo-indigo']
const placeholderClass = computed(() => placeholderPalette[(props.place.id ?? 0) % placeholderPalette.length])

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

const reviewers = computed(() => (props.place.reviewers || []).slice(0, 3))
</script>

<style scoped lang="scss">
.sb-result {
  display: flex;
  gap: 12px;
  padding: 12px 12px 14px;
  background: var(--sb-paper-card);
  border-radius: 1px;
  text-decoration: none;
  color: inherit;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.06),
    0 4px 14px rgba(40, 30, 20, 0.07);
  align-items: flex-start;
  position: relative;
  // Clip stray overflow from inner tilted polaroids when something extends past
  // the card on narrow viewports. Tape sits within thumb padding so it stays visible.
  overflow: hidden;
}

.thumb {
  flex-shrink: 0;
  position: relative;
  padding: 6px 6px 8px;
  // give room for tilt rotation and tape stub
}
.gem-corner {
  position: absolute;
  top: 4px;
  right: 4px;
  z-index: 2;
}

.body {
  flex: 1 1 0;
  min-width: 0;            // allow flex children to shrink past content size
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.title-row {
  display: flex;
  align-items: baseline;
  gap: 8px;
  min-width: 0;
}
.title {
  flex: 1 1 auto;
  min-width: 0;
  font-family: var(--sb-serif);
  font-style: italic;
  font-weight: 500;
  font-size: 17px;
  line-height: 1.15;
  color: var(--sb-ink);
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  word-break: break-word;
  hyphens: auto;
}
.gem-stamp {
  flex-shrink: 0;
}

.meta {
  font-family: var(--sb-hand);
  font-size: 15px;
  color: var(--sb-ink-mute);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.bottom {
  display: flex;
  align-items: center;
  gap: 10px 12px;
  margin-top: 6px;
  flex-wrap: wrap;
  min-width: 0;
}
// Mini ticket allows shrinking on very narrow screens
.mini-ticket {
  flex-shrink: 1;
}

.people {
  display: inline-flex;
  align-items: center;
  flex-shrink: 0;
}
// Stack avatars by negative margin — reset the absolute positioning baked in
// `.sb-author-tag` so they live in document flow inside .people.
.people .r-tag {
  position: relative;
  width: 22px;
  height: 22px;
  margin-left: -8px;
}
.people .r-tag:first-child {
  margin-left: 0;
}
.people .r-tag.has-photo {
  background: var(--sb-paper-card);
  overflow: hidden;
}
.people .r-ph {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
