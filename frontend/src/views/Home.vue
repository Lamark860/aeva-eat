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
            :class="[isFeatured(item, buckets.current.items) ? 'featured' : cellTilt(item)]"
          >
            <ArtifactCard
              v-if="item._kind !== 'note'"
              :place="item"
              :featured="isFeatured(item, buckets.current.items)"
            />
            <NoteArtifact
              v-else
              :note="item._note"
              :can-edit="auth.user?.id === item._note.author_id"
              @delete="onDeleteNote"
            />
          </div>
          <div
            v-if="item._kind !== 'note' && item.has_video"
            :key="`kr-${item.id}`"
            class="doska-cell divider-cell"
          >
            <KruzhokDivider :place-id="item.id" :place-name="item.name" :date="item.created_at" />
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

        <section v-else class="doska-archive">
          <div class="sb-week-head">
            <div class="dates">{{ b.label }}</div>
            <div class="rule"></div>
            <button class="sb-more-link collapse-btn" type="button" @click="toggleBucket(b.key)">
              свернуть&nbsp;↑
            </button>
          </div>

          <div class="doska-week">
            <template v-for="item in b.items" :key="`${b.key}-${item.id}`">
              <div
                class="doska-cell"
                :class="[isFeatured(item, b.items) ? 'featured' : cellTilt(item)]"
              >
                <ArtifactCard
                  v-if="item._kind !== 'note'"
                  :place="item"
                  :featured="isFeatured(item, b.items)"
                />
                <NoteArtifact
                  v-else
                  :note="item._note"
                  :can-edit="auth.user?.id === item._note.author_id"
                  @delete="onDeleteNote"
                />
              </div>
              <div
                v-if="item._kind !== 'note' && item.has_video"
                :key="`${b.key}-kr-${item.id}`"
                class="doska-cell divider-cell"
              >
                <KruzhokDivider :place-id="item.id" :place-name="item.name" :date="item.created_at" />
              </div>
            </template>
          </div>
        </section>
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
import KruzhokDivider from '../components/scrapbook/KruzhokDivider.vue'
import NoteArtifact from '../components/scrapbook/NoteArtifact.vue'
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

const currentVisible = computed(() => (buckets.value.current?.items || []).slice(0, visibleCount.value))

const tilts = ['sb-t-l3', 'sb-t-r1', 'sb-t-l2', 'sb-t-r2', 'sb-t-r3', 'sb-t-l1']
function cellTilt(item) {
  return tilts[(item.id ?? 0) % tilts.length]
}

// A1 — full-width «звезда» раз в 5+ карточек (NEXT.md §A1). Только для
// place-артефактов (записки никогда не становятся звездой). Приоритет:
// первая жемчужина > первый place-item (если в бакете ≥5 элементов).
function pickFeaturedId(items) {
  if (!items || items.length < 5) return null
  const placeItems = items.filter(i => i._kind !== 'note')
  if (!placeItems.length) return null
  const gem = placeItems.find(p => p.is_gem_place)
  return (gem ?? placeItems[0])?.id ?? null
}
function isFeatured(item, bucketItems) {
  if (item._kind === 'note') return false
  return item.id === pickFeaturedId(bucketItems)
}

function summaryFor(b) {
  return b.items.slice(0, 3).map((p, i) => ({
    id: p.id,
    src: p.image_url || '',
    cap: p.name,
    gem: !!p.is_gem_place,
    placeholder: ['sb-photo-warm', 'sb-photo-olive', 'sb-photo-dusk', 'sb-photo-sage', 'sb-photo-peach'][i % 5],
  }))
}
function gemsIn(b) {
  return b.items.filter((p) => p.is_gem_place).length
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

  &.featured {
    grid-column: 1 / -1;
    // У «звезды» убираем sb-t-* (утилитный tilt) — крупная карточка с
    // сильным наклоном выглядит навязчиво. Лёгкое покачивание оставляем.
    transform: rotate(-0.6deg);
  }

  // Kruzhok-разделитель — тоже full-width, без tilt'а, минимальная высота.
  &.divider-cell {
    grid-column: 1 / -1;
    transform: none;
  }
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
  color: #fff;
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
