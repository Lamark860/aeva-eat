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

    <!-- C1 — soft-CTA внизу Доски, только в onboarding-режиме -->
    <section v-if="onboarding" class="onboarding-cta">
      <div class="oc-hint">первое прикнопить —</div>
      <router-link
        :to="{ path: '/places/new', query: { intent: 'visit' } }"
        class="oc-button"
      >
        + новое место
      </router-link>
    </section>

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

// DESIGN-DECISIONS Q10: «бумага раскрывается» — height-animate 320мс на разворот
// прошедшей недели. Натуральная высота меряется через scrollHeight и анимируется
// в один проход; opacity stagger детей идёт через CSS animation-delay (см. .doska-archive .doska-cell).
const ARCHIVE_EASE = 'cubic-bezier(0.2, 0.8, 0.2, 1)'
function onArchiveEnter(el, done) {
  el.style.overflow = 'hidden'
  el.style.height = '0px'
  // force reflow so the transition starts from 0
  void el.offsetHeight
  el.style.transition = `height 320ms ${ARCHIVE_EASE}`
  el.style.height = el.scrollHeight + 'px'
  const onEnd = () => {
    el.style.height = ''
    el.style.transition = ''
    el.style.overflow = ''
    el.removeEventListener('transitionend', onEnd)
    done()
  }
  el.addEventListener('transitionend', onEnd)
}
function onArchiveLeave(el, done) {
  el.style.overflow = 'hidden'
  el.style.height = el.scrollHeight + 'px'
  void el.offsetHeight
  el.style.transition = `height 320ms ${ARCHIVE_EASE}`
  el.style.height = '0px'
  const onEnd = () => {
    el.removeEventListener('transitionend', onEnd)
    done()
  }
  el.addEventListener('transitionend', onEnd)
}

const currentVisible = computed(() => (buckets.value.current?.items || []).slice(0, visibleCount.value))

const tilts = ['sb-t-l3', 'sb-t-r1', 'sb-t-l2', 'sb-t-r2', 'sb-t-r3', 'sb-t-l1']
function cellTilt(item) {
  // id может быть числом (place) или строкой `note-N` / `wish-U-P`. Делаем
  // стабильный хэш строкового представления.
  const s = String(item.id ?? '0')
  let h = 0
  for (let i = 0; i < s.length; i++) h = (h * 31 + s.charCodeAt(i)) | 0
  return tilts[Math.abs(h) % tilts.length]
}

// A1 — full-width «звезда» раз в 5+ карточек (NEXT.md §A1). Только для
// place-артефактов (записки и wishlist-планы не становятся звездой).
// Приоритет: первая жемчужина > первый place-item (если в бакете ≥5 элементов).
function pickFeaturedId(items) {
  if (!items || items.length < 5) return null
  const placeItems = items.filter(i => i._kind === 'place')
  if (!placeItems.length) return null
  const gem = placeItems.find(p => p.is_gem_place)
  return (gem ?? placeItems[0])?.id ?? null
}
function isFeatured(item, bucketItems) {
  if (item._kind !== 'place') return false
  return item.id === pickFeaturedId(bucketItems)
}

function summaryFor(b) {
  // CollapsedStrip показывает превью только реальных визитов.
  return b.items.filter(i => i._kind === 'place').slice(0, 3).map((p, i) => ({
    id: p.id,
    src: p.image_url || '',
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
  // dense — мелкие артефакты подтягиваются под full-width к свободным
  // местам, чтобы после «звезды» не оставался зазор справа.
  grid-auto-flow: dense;
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

  // Q-video: место с видео получает всю строку. Кружочек прилипает к карточке
  // изнутри ArtifactCard — отдельной KruzhokDivider-ячейки больше нет.
  // Дополнительный padding-bottom — чтобы тень кружочка не наезжала на следующий ряд.
  &.has-video {
    grid-column: 1 / -1;
    padding-bottom: 14px;
  }
}

@keyframes sb-cell-in {
  from { opacity: 0; translate: 0 8px; }
  to   { opacity: 1; translate: 0 0; }
}

// Q10: при разворачивании архива дети вступают со stagger 40мс — «бумага раскрывается».
.archive-stagger .doska-cell {
  animation-delay: calc(var(--i, 0) * 40ms);
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

.onboarding-cta {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 22px 18px 28px;
}
.oc-hint {
  font-family: var(--sb-hand);
  font-size: 17px;
  color: var(--sb-ink-mute);
}
.oc-button {
  display: inline-flex;
  align-items: center;
  padding: 12px 22px;
  background: var(--sb-terracotta);
  color: var(--sb-on-accent);
  border-radius: 999px;
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 16px;
  text-decoration: none;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.1),
    0 4px 10px rgba(40, 30, 20, 0.15);
  transition: transform 200ms ease;
  &:hover { transform: translateY(-1px); }
}
</style>
