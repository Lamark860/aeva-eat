<template>
  <div class="sb-paper sb-grain sb-screen doska">
    <header class="doska-header">
      <div class="sb-wordmark">aeva<span class="dot"></span>eat</div>
      <div class="sub">наша доска</div>
    </header>

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
        <div class="col">
          <div
            v-for="item in colA"
            :key="`a-${item.id}`"
            class="doska-cell"
            :class="cellTilt(item, 'a')"
          >
            <ArtifactCard :place="item" />
          </div>
        </div>
        <div class="col">
          <div
            v-for="item in colB"
            :key="`b-${item.id}`"
            class="doska-cell"
            :class="cellTilt(item, 'b')"
          >
            <ArtifactCard :place="item" />
          </div>
        </div>
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
            <div class="col">
              <div
                v-for="item in bucketColA(b)"
                :key="`${b.key}-a-${item.id}`"
                class="doska-cell"
                :class="cellTilt(item, 'a')"
              >
                <ArtifactCard :place="item" />
              </div>
            </div>
            <div class="col">
              <div
                v-for="item in bucketColB(b)"
                :key="`${b.key}-b-${item.id}`"
                class="doska-cell"
                :class="cellTilt(item, 'b')"
              >
                <ArtifactCard :place="item" />
              </div>
            </div>
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

    <AddArtifactSheet v-model:open="openSheet" @pick="onPick" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useFeed } from '../composables/useFeed'
import PinButton from '../components/scrapbook/PinButton.vue'
import CollapsedStrip from '../components/scrapbook/CollapsedStrip.vue'
import AddArtifactSheet from '../components/scrapbook/AddArtifactSheet.vue'
import ArtifactCard from '../components/scrapbook/ArtifactCard.vue'

const router = useRouter()
const { places, loading, buckets, load } = useFeed()

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

function bucketColA(b) { return b.items.filter((_, i) => i % 2 === 0) }
function bucketColB(b) { return b.items.filter((_, i) => i % 2 === 1) }

const currentVisible = computed(() => (buckets.value.current?.items || []).slice(0, visibleCount.value))

// Split visible items into two masonry-style columns
const colA = computed(() => currentVisible.value.filter((_, i) => i % 2 === 0))
const colB = computed(() => currentVisible.value.filter((_, i) => i % 2 === 1))

const tiltsA = ['sb-t-l3', 'sb-t-r1', 'sb-t-l2', 'sb-t-r2']
const tiltsB = ['sb-t-r2', 'sb-t-l2', 'sb-t-r3', 'sb-t-l1']
function cellTilt(item, col) {
  const arr = col === 'a' ? tiltsA : tiltsB
  const idx = (item.id ?? 0) % arr.length
  return arr[idx]
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
  if (kind === 'visit') router.push({ path: '/places/new', query: { intent: 'visit' } })
}

onMounted(load)
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
  display: flex;
  gap: 14px;
  // Extra bottom padding so tilted polaroid shadows don't visually collide
  // with whatever follows (collapsed strip / archive footer).
  padding: 14px 16px 26px;
  align-items: flex-start;

  .col {
    flex: 1 1 0;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 28px;
  }
}

.doska-cell {
  // tilts come from sb-t-* utility classes — leave room for them
  transform-origin: center;
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
</style>
