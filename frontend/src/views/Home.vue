<template>
  <div class="sb-paper sb-grain sb-screen doska">
    <header class="doska-header">
      <div class="sb-wordmark">aeva<span class="dot"></span>eat</div>
      <div class="sub">наша доска</div>
    </header>

    <!-- C1 — onboarding для нового пригласённого. Видит «обзорный режим»:
         история круга read-only + рукописный приветственный хинт сверху.
         Режим выключается после первого review (places_count > 0). -->
    <section v-if="onboarding" class="onboarding-banner">
      <div class="ob-line">здесь живёт камерный дневник еды круга</div>
      <div class="ob-sub">пролистай — а когда захочется, прикнопь первое место</div>
    </section>

    <section v-if="loading && !places.length" class="sb-empty">
      листаем заметки…
    </section>

    <template v-else>
      <!-- Current week -->
      <div v-if="buckets.current" class="sb-week-head">
        <div class="dates">{{ buckets.current.label }}</div>
        <div class="rule"></div>
        <PinButton @click="openSheet = true">добавить</PinButton>
      </div>

      <div v-if="currentVisible.length" class="doska-week">
        <template v-for="item in currentVisible" :key="item.id">
          <div
            class="doska-cell"
            :class="[
              cellTilt(item),
              {
                featured: isFeatured(item, buckets.current.items),
                'has-video': item._kind === 'place' && !!(item.videos?.length || item.video_url),
              },
            ]"
          >
            <ArtifactCard
              v-if="item._kind === 'place'"
              :place="item"
              :featured="isFeatured(item, buckets.current.items)"
              :attendees="item._attendees"
            />
            <NoteArtifact
              v-else-if="item._kind === 'note'"
              :note="item._note"
              :can-edit="auth.user?.id === item._note.author_id"
              @delete="onDeleteNote"
            />
            <WishlistArtifact
              v-else-if="item._kind === 'wishlist'"
              :entry="item._wish"
            />
          </div>
        </template>
      </div>

      <div v-else class="sb-empty">
        <div>пусто на этой неделе</div>
        <div style="margin-top: 8px; font-size: 16px">
          прикнопь первое место —
        </div>
        <div style="margin-top: 12px">
          <PinButton @click="openSheet = true">добавить</PinButton>
        </div>
      </div>

      <div
        v-if="buckets.current && buckets.current.items.length > visibleCount"
        style="text-align: center; padding: 4px 0 4px"
      >
        <button class="sb-more-link" type="button" @click="visibleCount += 6">
          ↓ ещё {{ buckets.current.items.length - visibleCount }} на этой неделе
        </button>
      </div>

      <!-- Older buckets — collapsed strip OR expanded inline feed -->
      <template v-for="b in buckets.older" :key="b.key">
        <CollapsedStrip
          v-if="!isExpanded(b.key)"
          :dates="b.label"
          :summary="summaryFor(b)"
          :count="b.items.length"
          :gem-count="gemsIn(b)"
          @expand="toggleBucket(b.key)"
        />

        <Transition name="archive" @enter="onArchiveEnter" @leave="onArchiveLeave">
          <section v-if="isExpanded(b.key)" class="doska-archive archive-stagger">
            <div class="sb-week-head">
              <div class="dates">{{ b.label }}</div>
              <div class="rule"></div>
              <button class="sb-more-link collapse-btn" type="button" @click="toggleBucket(b.key)">
                свернуть&nbsp;↑
              </button>
            </div>

            <div class="doska-week">
              <template v-for="(item, i) in b.items" :key="`${b.key}-${item.id}`">
                <div
                  class="doska-cell"
                  :class="[
                    cellTilt(item),
                    {
                      featured: isFeatured(item, b.items),
                      'has-video': item._kind === 'place' && !!(item.videos?.length || item.video_url),
                    },
                  ]"
                  :style="{ '--i': i }"
                >
                  <ArtifactCard
                    v-if="item._kind === 'place'"
                    :place="item"
                    :featured="isFeatured(item, b.items)"
                    :attendees="item._attendees"
                  />
                  <NoteArtifact
                    v-else-if="item._kind === 'note'"
                    :note="item._note"
                    :can-edit="auth.user?.id === item._note.author_id"
                    @delete="onDeleteNote"
                  />
                  <WishlistArtifact
                    v-else-if="item._kind === 'wishlist'"
                    :entry="item._wish"
                  />
                </div>
              </template>
            </div>
          </section>
        </Transition>
      </template>

      <div
        v-if="buckets.older.length"
        style="text-align: center; padding: 12px 0 6px; font-family: var(--sb-hand); font-size: 16px; color: var(--sb-ink-mute)"
      >
        ↓&nbsp; раньше — в архиве
      </div>

      <div
        v-else-if="!loading && places.length && !buckets.current?.items.length"
        style="text-align: center; padding: 24px 0; font-family: var(--sb-hand); font-size: 16px; color: var(--sb-ink-mute)"
      >
        пока всё, что есть
      </div>
    </template>

    <!-- D5 — onboarding-CTA внизу убран. Один призыв (баннер сверху) сильнее
         двух; на пустой неделе уже виден большой PinButton «добавить» — этого
         достаточно. См. R6 D5. -->

    <AddArtifactSheet v-model:open="openSheet" @pick="onPick" />
    <NoteSheet v-model:open="noteFormOpen" @submit="onSubmitNote" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useFeed } from '../composables/useFeed'
import PinButton from '../components/scrapbook/PinButton.vue'
import CollapsedStrip from '../components/scrapbook/CollapsedStrip.vue'
import AddArtifactSheet from '../components/scrapbook/AddArtifactSheet.vue'
import NoteSheet from '../components/scrapbook/NoteSheet.vue'
import ArtifactCard from '../components/scrapbook/ArtifactCard.vue'
import NoteArtifact from '../components/scrapbook/NoteArtifact.vue'
import WishlistArtifact from '../components/scrapbook/WishlistArtifact.vue'
import { useAuthStore } from '../stores/auth'
import { useNotesStore } from '../stores/notes'
import { useToast } from '../composables/useToast'
import http from '../api/http'

const router = useRouter()
const auth = useAuthStore()
const notesStore = useNotesStore()
const toast = useToast()
const { places, loading, buckets, load } = useFeed()

// C1 — определяем onboarding-режим: пользователь без визитов. Подгружаем
// его публичный профиль (place_count) единожды; ленту всё равно показываем
// — это и есть «обзорный режим» (read-only история круга).
const myProfile = ref(null)
async function loadMyProfile() {
  if (!auth.user?.id) return
  try {
    const { data } = await http.get(`/users/${auth.user.id}`)
    myProfile.value = data
  } catch { /* ignore — индикатор без профиля просто не сработает */ }
}
const onboarding = computed(() =>
  myProfile.value !== null && (myProfile.value.place_count || 0) === 0
)

const openSheet = ref(false)
const visibleCount = ref(4)
const expandedBuckets = ref(new Set())

function isExpanded(key) { return expandedBuckets.value.has(key) }
function toggleBucket(key) {
  const next = new Set(expandedBuckets.value)
  if (next.has(key)) next.delete(key)
  else next.add(key)
  expandedBuckets.value = next
}

// DESIGN-DECISIONS Q10: «бумага раскрывается» — высота + прозрачность анимируются
// вместе (рост без fade выглядит как жёсткий клип). Натуральная высота меряется
// через scrollHeight; дети въезжают CSS-staggers'ом (см. .archive-stagger .doska-cell).
// transitionend слушаем только по height; setTimeout-fallback на случай если он не стрельнёт.
const ARCHIVE_EASE = 'cubic-bezier(0.22, 0.61, 0.36, 1)'
const ARCHIVE_MS = 360
function animateHeight(el, from, to, fade, done) {
  el.style.overflow = 'hidden'
  el.style.willChange = 'height, opacity'
  el.style.height = from
  el.style.opacity = fade[0]
  // force reflow so the transition starts from `from`
  void el.offsetHeight
  el.style.transition = `height ${ARCHIVE_MS}ms ${ARCHIVE_EASE}, opacity ${ARCHIVE_MS - 100}ms ease`
  el.style.height = to
  el.style.opacity = fade[1]
  let finished = false
  const cleanup = (e) => {
    if (e && e.propertyName !== 'height') return
    if (finished) return
    finished = true
    el.style.height = ''
    el.style.transition = ''
    el.style.overflow = ''
    el.style.opacity = ''
    el.style.willChange = ''
    el.removeEventListener('transitionend', cleanup)
    done()
  }
  el.addEventListener('transitionend', cleanup)
  setTimeout(cleanup, ARCHIVE_MS + 90)
}
function onArchiveEnter(el, done) {
  animateHeight(el, '0px', el.scrollHeight + 'px', ['0', '1'], done)
}
function onArchiveLeave(el, done) {
  animateHeight(el, el.scrollHeight + 'px', '0px', ['1', '0'], done)
}

const currentVisible = computed(() => (buckets.value.current?.items || []).slice(0, visibleCount.value))

const tilts = ['sb-t-l3', 'sb-t-r1', 'sb-t-l2', 'sb-t-r2', 'sb-t-r3', 'sb-t-l1']
// R5-Q4 — wishlist'у даём только мягкие tilt'ы (±1° или ±2°). Резкий ±3°
// на «плане» читается агрессивно, как живое впечатление, чего wishlist
// концептуально не несёт.
const softTilts = ['sb-t-l1', 'sb-t-r1', 'sb-t-l2', 'sb-t-r2']
function cellTilt(item) {
  // id может быть числом (place) или строкой `note-N` / `wish-U-P`. Делаем
  // стабильный хэш строкового представления.
  const s = String(item.id ?? '0')
  let h = 0
  for (let i = 0; i < s.length; i++) h = (h * 31 + s.charCodeAt(i)) | 0
  const pool = item._kind === 'wishlist' ? softTilts : tilts
  return pool[Math.abs(h) % pool.length]
}

// D2 — ровно ОДНА full-width «звезда» на неделю. Иерархия выбора:
// 1) первая жемчужина (статусный артефакт сильнее),
// 2) первый place с видео-кружком (кружок крупнее на full-width),
// 3) первый place любой природы.
// Требование на ≥5 элементов сохранили — на тонких неделях звезда лишняя.
// Записки и wishlist-планы звездой не становятся (концептуально не «находка»).
function pickFeaturedId(items) {
  if (!items || items.length < 5) return null
  const placeItems = items.filter(i => i._kind === 'place')
  if (!placeItems.length) return null
  const gem = placeItems.find(p => p.is_gem_place)
  if (gem) return gem.id
  const withVideo = placeItems.find(p => (p.video_url || (p.videos || []).length > 0))
  if (withVideo) return withVideo.id
  return placeItems[0]?.id ?? null
}
function isFeatured(item, bucketItems) {
  if (item._kind !== 'place') return false
  return item.id === pickFeaturedId(bucketItems)
}

function summaryFor(b) {
  // CollapsedStrip показывает превью только реальных визитов.
  // Приоритезируем места с фото (cover или review-фото) — иначе полоса
  // на 8 мест с одним фото выглядит как 3 пустых полароида.
  const places = b.items.filter(i => i._kind === 'place')
  const coverOf = (p) => p.image_url || p.feed_photos?.[0]?.url || ''
  const withPhoto = places.filter(p => coverOf(p))
  const withoutPhoto = places.filter(p => !coverOf(p))
  const ordered = [...withPhoto, ...withoutPhoto].slice(0, 3)
  return ordered.map((p, i) => ({
    id: p.id,
    src: coverOf(p),
    cap: p.name,
    gem: !!p.is_gem_place,
    placeholder: ['sb-photo-warm', 'sb-photo-olive', 'sb-photo-dusk', 'sb-photo-sage', 'sb-photo-peach'][i % 5],
  }))
}
function gemsIn(b) {
  return b.items.filter((p) => p._kind === 'place' && p.is_gem_place).length
}

function onPick(kind) {
  if (kind === 'visit') {
    router.push({ path: '/places/new', query: { intent: 'visit' } })
  } else if (kind === 'note') {
    noteFormOpen.value = true
  }
}

const noteFormOpen = ref(false)
async function onSubmitNote(payload) {
  try {
    await notesStore.create(payload)
    toast.success('Записка прикноплена')
    await load()
  } catch (e) {
    toast.error(e.response?.data?.error || 'Не удалось сохранить записку')
  }
  noteFormOpen.value = false
}

async function onDeleteNote(id) {
  if (!confirm('Удалить записку?')) return
  try {
    await notesStore.remove(id)
    await load()
    toast.info('Записка убрана')
  } catch (e) {
    toast.error(e.response?.data?.error || 'Не удалось удалить')
  }
}

onMounted(() => {
  load()
  loadMyProfile()
})
</script>

<style scoped lang="scss">
.doska {
  // Wordmark + sub line at the very top, with safe-area inset.
  padding-top: calc(18px + var(--aeva-safe-top, 0px));
}

.doska-header {
  padding: 0 18px 4px;
  display: flex;
  align-items: baseline;
  gap: 10px;

  .sub {
    font-family: var(--sb-hand);
    font-size: 18px;
    color: var(--sb-ink-mute);
    line-height: 1;
  }
}

.doska-week {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  // D3 — `dense` убран: он подтягивает мелкие артефакты к свободным
  // ячейкам и тем самым ломает хронологию (записка вторника визуально
  // оказывается под визитом среды). Пустая ячейка после full-width
  // считается приемлемой ценой за читаемую последовательность.
  gap: 28px 14px;
  // Extra bottom padding so tilted polaroid shadows don't visually collide
  // with whatever follows (collapsed strip / archive footer).
  padding: 14px 16px 26px;
}

.doska-cell {
  // tilts come from sb-t-* utility classes — leave room for them
  transform-origin: center;
  min-width: 0;
  // DESIGN-DECISIONS Q10: первый mount артефакта — fade + slide-up 8px, 280мс.
  // `translate` (а не `transform`), чтобы не клобберить tilt-rotate из sb-t-*.
  animation: sb-cell-in 280ms cubic-bezier(0.2, 0.8, 0.2, 1) both;

  &.featured {
    grid-column: 1 / -1;
    // У «звезды» убираем sb-t-* (утилитный tilt) — крупная карточка с
    // сильным наклоном выглядит навязчиво. Лёгкое покачивание оставляем.
    transform: rotate(-0.6deg);
  }

  // Q-video: чтобы тень видео-кружочка не наезжала на следующий ряд.
  // Full-width больше не выдаём автоматически — это решается через .featured
  // (см. pickFeaturedId, иерархия gem > video > first).
  &.has-video {
    padding-bottom: 14px;
  }
}

@keyframes sb-cell-in {
  from { opacity: 0; translate: 0 8px; }
  to   { opacity: 1; translate: 0 0; }
}

// Q10: при разворачивании архива дети вступают со stagger — «бумага раскрывается».
// Индекс ограничиваем (min ...), иначе на длинной неделе хвост докручивается
// заметно дольше, чем растёт высота контейнера, и это читается как лаг.
.archive-stagger .doska-cell {
  animation-delay: calc(min(var(--i, 0), 5) * 30ms);
}

// Доступность: при системном «уменьшить движение» гасим mount-анимации
// (контейнер всё равно раскроется по высоте — это не вестибулярная анимация).
@media (prefers-reduced-motion: reduce) {
  .doska-cell { animation: none; }
}

.doska-archive {
  margin-top: 8px;
  // Same buffer as the current-week container — keep adjacent buckets clear of tilt overlap.
  padding-bottom: 8px;
  // DESIGN-DECISIONS L4: чуть приглушённая бумага, чтобы развёрнутый архив
  // визуально читался как «прошлое». Деликатно, ~5–7% разницы.
  background: var(--sb-paper-deep);
}

.collapse-btn {
  margin-left: auto;
  padding: 4px 8px;
}

/* C1 onboarding — приветственный хинт сверху + soft-CTA снизу. */
.onboarding-banner {
  text-align: center;
  padding: 8px 18px 18px;
}
.ob-line {
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 17px;
  color: var(--sb-ink);
  line-height: 1.3;
}
.ob-sub {
  font-family: var(--sb-hand);
  font-size: 17px;
  color: var(--sb-ink-mute);
  margin-top: 4px;
}

</style>
