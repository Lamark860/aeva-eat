<template>
  <div class="sb-paper sb-grain sb-screen detail">
    <div v-if="placesStore.loading && !place" class="sb-empty">
      листаем заметки…
    </div>

    <template v-else-if="place">
      <!-- Top: cover polaroid + title block -->
      <section class="detail-top">
        <div class="cover">
          <Polaroid
            :src="place.image_url || ''"
            :width="280"
            :height="220"
            :tilt="coverTilt"
            :gem="!!place.is_gem_place"
            :placeholder="coverPlaceholder"
          >
            <Tape :variant="coverTapeVariant" :style="coverTapeStyle" />
            <span v-if="place.is_gem_place" class="cover-gem">
              <GemBadge :size="28" />
            </span>
          </Polaroid>
        </div>

        <h1 class="title">{{ place.name }}</h1>

        <div class="stamps">
          <Stamp v-if="place.city" kind="ink">{{ place.city }}</Stamp>
          <Stamp v-if="place.cuisine_type">{{ place.cuisine_type }}</Stamp>
          <Stamp
            v-for="cat in place.categories || []"
            :key="cat"
            kind="moss"
          >
{{ cat }}
</Stamp>
          <Stamp v-if="place.is_gem_place" kind="gem">жемчужина</Stamp>
        </div>

        <div v-if="place.address" class="hand-meta">{{ place.address }}</div>
        <div v-if="place.website" class="hand-meta">
          <a :href="place.website" target="_blank" rel="noopener">{{ websiteHost }}</a>
        </div>
      </section>

      <!-- Summary ratings -->
      <section v-if="place.review_count > 0" class="detail-summary">
        <div class="overall">
          <div class="num">{{ overallRating }}</div>
          <div class="lbl">общая</div>
        </div>
        <div class="ticket-wrap sb-t-l1">
          <Ticket
            :food="place.avg_food"
            :service="place.avg_service"
            :vibe="place.avg_vibe"
          />
        </div>
        <div class="reviews-count">
          <div class="num">{{ place.review_count }}</div>
          <div class="lbl">{{ reviewsLabel }}</div>
        </div>
      </section>

      <!-- CTA buttons -->
      <section class="detail-cta">
        <button
          v-if="auth.isAuthenticated && !isOwner"
          class="cta-pin"
          :class="{ active: wishlist.isWishlisted(place.id) }"
          @click="wishlist.toggle(place.id)"
        >
          <span class="head" aria-hidden="true"></span>
          <span class="lbl">{{ wishlist.isWishlisted(place.id) ? 'в&nbsp;планах' : 'в&nbsp;wishlist' }}</span>
        </button>

        <a v-if="hasMap" href="#detail-map" class="cta-link">↗ показать на карте</a>

        <button
          v-if="auth.isAuthenticated && !showForm && !editingReview"
          type="button"
          class="cta-link"
          @click="openForm"
        >
          ✎ оставить отзыв
        </button>

        <template v-if="isOwner">
          <router-link :to="`/places/${place.id}/edit`" class="cta-link">ред. место</router-link>
          <button class="cta-link danger" @click="handleDelete">удалить</button>
        </template>
      </section>

      <!-- Mini map -->
      <section v-if="hasMap" id="detail-map" class="detail-map">
        <MapView
          :places="[place]"
          :center="[place.lat, place.lng]"
          :zoom="15"
          height="240px"
          :single-marker="true"
        />
      </section>

      <div class="sb-section-head" style="padding: 18px 18px 8px">
        <h2>Отзывы</h2>
        <span class="sub">кто что прикнопил</span>
      </div>

      <!-- Existing reviews come first — what others wrote about this place. -->
      <section class="detail-reviews">
        <div v-if="reviewsStore.loading" class="sb-empty">…</div>
        <div v-else-if="reviewsStore.reviews.length === 0" class="sb-empty">
          ничего ещё не прикнопили — оставь первое
        </div>
        <ReviewCard
          v-for="rv in reviewsStore.reviews"
          :key="rv.id"
          :review="rv"
          :can-edit="canEditReview(rv)"
          @edit="onEdit(rv)"
          @delete="handleDeleteReview"
        />
      </section>

      <!-- Review form — collapsed by default, toggled by the «оставить отзыв» CTA above
           or shown automatically when editing an existing review. -->
      <section
        v-if="(showForm && auth.isAuthenticated) || editingReview"
        id="detail-review-form"
        class="detail-form"
      >
        <ReviewForm
          v-if="!editingReview"
          :key="reviewFormKey"
          :place-id="place.id"
          @submitted="handleCreateReview"
        />
        <ReviewForm
          v-else
          :place-id="place.id"
          :review="editingReview"
          @submitted="handleUpdateReview"
          @cancel="editingReview = null"
        />
      </section>

      <div class="detail-back">
        <router-link to="/places" class="cta-link">← назад к списку</router-link>
      </div>
    </template>

    <div v-else class="sb-empty">место не найдено</div>
  </div>
</template>

<script setup>
import { onMounted, computed, ref, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePlacesStore } from '../stores/places'
import { useReviewsStore } from '../stores/reviews'
import { useAuthStore } from '../stores/auth'
import { useWishlistStore } from '../stores/wishlist'
import { useToast } from '../composables/useToast'
import ReviewCard from '../components/ReviewCard.vue'
import ReviewForm from '../components/ReviewForm.vue'
import MapView from '../components/MapView.vue'
import Polaroid from '../components/scrapbook/Polaroid.vue'
import Tape from '../components/scrapbook/Tape.vue'
import Stamp from '../components/scrapbook/Stamp.vue'
import Ticket from '../components/scrapbook/Ticket.vue'
import GemBadge from '../components/scrapbook/GemBadge.vue'

const route = useRoute()
const router = useRouter()
const placesStore = usePlacesStore()
const reviewsStore = useReviewsStore()
const auth = useAuthStore()
const wishlist = useWishlistStore()
const toast = useToast()

const editingReview = ref(null)
const reviewFormKey = ref(0)
const showForm = ref(false)

function openForm() {
  showForm.value = true
  // wait for the form to render before scrolling to it
  nextTick(() => {
    document.getElementById('detail-review-form')?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  })
}

function onEdit(rv) {
  editingReview.value = rv
  showForm.value = false
  nextTick(() => {
    document.getElementById('detail-review-form')?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  })
}

const place = computed(() => placesStore.currentPlace)
const isOwner = computed(() => auth.user && place.value && place.value.created_by === auth.user.id)
const hasMap = computed(() => place.value && place.value.lat && place.value.lng)

const overallRating = computed(() => {
  const p = place.value
  if (!p || p.avg_food == null || p.avg_service == null || p.avg_vibe == null) return '–'
  return ((Number(p.avg_food) + Number(p.avg_service) + Number(p.avg_vibe)) / 3).toFixed(1)
})

const reviewsLabel = computed(() => {
  const n = (place.value?.review_count || 0) % 100
  if (n >= 11 && n <= 14) return 'отзывов'
  const last = n % 10
  if (last === 1) return 'отзыв'
  if (last >= 2 && last <= 4) return 'отзыва'
  return 'отзывов'
})

const websiteHost = computed(() => {
  try { return new URL(place.value.website).host } catch { return place.value.website }
})

// Cover styling — deterministic per place id.
const tilts = ['t-l3', 't-r2', 't-l2']
const coverTilt = computed(() => tilts[(place.value?.id ?? 0) % tilts.length])

const coverTapeVariant = computed(() => ['', 'rose', 'mint', 'blue'][(place.value?.id ?? 0) % 4])
const coverTapeStyle = computed(() => {
  const variants = [
    { top: '-12px', left: '50%', transform: 'translateX(-50%) rotate(-8deg)', width: '92px' },
    { top: '-10px', left: '32px', transform: 'rotate(-12deg)' },
    { top: '-10px', right: '32px', transform: 'rotate(10deg)' },
  ]
  return variants[(place.value?.id ?? 0) % variants.length]
})

const placeholderPalette = ['sb-photo-warm', 'sb-photo-olive', 'sb-photo-dusk', 'sb-photo-sage', 'sb-photo-peach', 'sb-photo-brick', 'sb-photo-cream', 'sb-photo-slate', 'sb-photo-indigo']
const coverPlaceholder = computed(() => placeholderPalette[(place.value?.id ?? 0) % placeholderPalette.length])

function canEditReview(rv) {
  if (!auth.user) return false
  return rv.authors?.some((a) => a.id === auth.user.id)
}

async function handleCreateReview(data) {
  const photoFile = data._photoFile
  const videoFile = data._videoFile
  delete data._photoFile
  delete data._videoFile
  try {
    const created = await reviewsStore.createReview(place.value.id, data)
    if (photoFile && created?.id) {
      await reviewsStore.uploadReviewImage(place.value.id, created.id, photoFile)
    }
    if (videoFile && created?.id) {
      await reviewsStore.uploadReviewVideo(place.value.id, created.id, videoFile)
    }
    toast.success('Отзыв добавлен!')
  } catch (e) {
    toast.error(e.response?.data?.error || 'Ошибка при создании отзыва')
  }
  reviewFormKey.value++
  showForm.value = false
  await placesStore.fetchPlace(route.params.id)
  await reviewsStore.fetchByPlace(route.params.id)
}

async function handleUpdateReview(data) {
  const photoFile = data._photoFile
  const videoFile = data._videoFile
  delete data._photoFile
  delete data._videoFile
  try {
    await reviewsStore.updateReview(place.value.id, editingReview.value.id, data)
    if (photoFile) {
      await reviewsStore.uploadReviewImage(place.value.id, editingReview.value.id, photoFile)
    }
    if (videoFile) {
      await reviewsStore.uploadReviewVideo(place.value.id, editingReview.value.id, videoFile)
    }
    editingReview.value = null
    toast.success('Отзыв обновлён')
  } catch (e) {
    toast.error(e.response?.data?.error || 'Ошибка при обновлении отзыва')
  }
  await placesStore.fetchPlace(route.params.id)
  await reviewsStore.fetchByPlace(route.params.id)
}

async function handleDeleteReview(id) {
  if (!confirm('Удалить отзыв?')) return
  await reviewsStore.deleteReview(place.value.id, id)
  toast.info('Отзыв удалён')
  await placesStore.fetchPlace(route.params.id)
}

async function handleDelete() {
  if (!confirm('Удалить заведение?')) return
  await placesStore.deletePlace(place.value.id)
  toast.info('Заведение удалено')
  router.push('/places')
}

onMounted(async () => {
  await placesStore.fetchPlace(route.params.id)
  await reviewsStore.fetchByPlace(route.params.id)

  // Coming from "+ добавить → новый визит" → /places/new?intent=visit redirected
  // here with ?review=open. Honor it: open the form, scroll to it, then drop
  // the query so a refresh doesn't re-trigger.
  if (route.query.review === 'open' && auth.isAuthenticated) {
    showForm.value = true
    await nextTick()
    document.getElementById('detail-review-form')?.scrollIntoView({ behavior: 'smooth', block: 'start' })
    router.replace({ path: route.path, query: {} })
  }
})
</script>

<style scoped lang="scss">
.detail {
  // Keep `.sb-screen` bottom padding (reserves space for the fixed BottomTabBar)
  // — only override top + sides here.
  padding-top: calc(18px + var(--aeva-safe-top, 0px));
  padding-inline: 0;
}

.detail-top {
  padding: 0 18px 8px;
  text-align: center;
}
.cover {
  display: flex;
  justify-content: center;
  margin: 0 0 18px;
  position: relative;
}
.cover-gem {
  position: absolute;
  top: 14px;
  right: 14px;
  z-index: 2;
}

.title {
  font-family: var(--sb-serif);
  font-style: italic;
  font-weight: 500;
  font-size: 30px;
  line-height: 1.1;
  color: var(--sb-ink);
  margin: 0 0 10px;
  word-break: break-word;
  overflow-wrap: break-word;
  hyphens: auto;
}

.stamps {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  justify-content: center;
  margin: 0 0 8px;
}

.hand-meta {
  font-family: var(--sb-hand);
  font-size: 16px;
  color: var(--sb-ink-mute);
  margin-top: 4px;

  a { color: inherit; text-decoration: underline; }
}

.detail-summary {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 18px;
  padding: 14px 18px 6px;
  flex-wrap: wrap;

  .overall, .reviews-count {
    text-align: center;

    .num {
      font-family: var(--sb-serif);
      font-style: italic;
      font-weight: 500;
      font-size: 32px;
      color: var(--sb-terracotta);
      line-height: 1;
    }
    .lbl {
      font-family: var(--sb-hand);
      font-size: 14px;
      color: var(--sb-ink-mute);
      margin-top: 2px;
    }
  }
  .reviews-count .num { color: var(--sb-ink); }
}

.detail-cta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 14px;
  justify-content: center;
  align-items: center;
  padding: 8px 16px 14px;
}

.cta-pin {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 14px 8px 10px;
  background: oklch(0.93 0.04 85);
  color: var(--sb-ink);
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
  border-radius: 999px;
  border: none;
  cursor: pointer;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.1),
    0 3px 8px rgba(40, 30, 20, 0.14);
  min-height: 36px;

  .head {
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background: radial-gradient(circle at 35% 30%, oklch(0.7 0.16 25), oklch(0.42 0.18 25));
    box-shadow:
      inset 0 -1px 1px rgba(0, 0, 0, 0.3),
      inset 0 1px 1px rgba(255, 255, 255, 0.3),
      0 1px 2px rgba(40, 15, 5, 0.3);
  }

  &.active {
    background: oklch(0.92 0.07 25);
    color: var(--sb-terracotta);
  }
}

.cta-link {
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
  color: var(--sb-ink);
  text-decoration: none;
  background: transparent;
  border: none;
  padding: 6px 4px;
  cursor: pointer;
  &:hover { color: var(--sb-terracotta); }
  &.danger:hover { color: var(--sb-terracotta); }
}

.detail-map {
  margin: 0 16px 18px;
  border-radius: 4px;
  overflow: hidden;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.06),
    0 4px 14px rgba(40, 30, 20, 0.07);
}

.detail-form {
  margin: 0 16px 18px;
  padding: 14px;
  background: #fdfcf7;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.06),
    0 4px 14px rgba(40, 30, 20, 0.07);
  border-radius: 1px;
}

.detail-reviews {
  padding: 0 16px;
}

.detail-back {
  text-align: center;
  padding: 16px 0 12px;
}

.ticket-wrap { display: inline-block; }
</style>
