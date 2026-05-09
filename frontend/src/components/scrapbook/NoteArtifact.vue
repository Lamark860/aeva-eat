<template>
  <!--
    Записка на доске. Бумага на тейпе, рукописный текст, опционально —
    привязка к месту/городу через хвостик-ссылку. Используется в ленте
    Home.vue рядом с ArtifactCard (B2 frontend).
  -->
  <article class="sb-note-art" :class="{ struck: note.is_struck }">
    <Note lined>
      <header class="na-head">
        <span
          class="r-tag sb-author-tag"
          :class="[authorColor(note.author_id), { 'has-photo': hasAvatar }]"
          :title="authorName"
        >
          <img v-if="hasAvatar" :src="note.author?.avatar_url" alt="" class="r-ph" />
          <template v-else>{{ authorInitial }}</template>
        </span>
        <span class="na-name">{{ authorName }}</span>
        <span class="na-date">{{ dateLabel }}</span>
      </header>

      <p class="na-text" :class="{ struck: note.is_struck }">{{ note.text }}</p>

      <footer v-if="note.place_id || note.city" class="na-foot">
        <router-link
          v-if="note.place_id"
          :to="`/places/${note.place_id}`"
          class="na-link"
        >
          → {{ note.place_name || `место #${note.place_id}` }}
        </router-link>
        <span v-else class="na-city">в&nbsp;{{ note.city }}</span>
      </footer>
    </Note>

    <Tape :variant="tapeVariant" :style="tapeStyle" />

    <button
      v-if="canEdit"
      class="na-x"
      type="button"
      :aria-label="'удалить'"
      @click="$emit('delete', note.id)"
    >
      ×
    </button>
  </article>
</template>

<script setup>
import { computed } from 'vue'
import Note from './Note.vue'
import Tape from './Tape.vue'
import { authorColor } from '../../composables/useFeed'

const props = defineProps({
  note:    { type: Object, required: true },
  canEdit: { type: Boolean, default: false },
})
defineEmits(['delete'])

const authorName = computed(() => props.note.author?.username || `автор #${props.note.author_id}`)
const authorInitial = computed(() => (authorName.value || '?').slice(0, 1).toUpperCase())
const hasAvatar = computed(() => !!props.note.author?.avatar_url)

const tapeVariants = ['', 'rose', 'mint', 'blue']
const tapeVariant = computed(() => {
  const palette = props.note.tape_color
  if (palette && tapeVariants.includes(palette)) return palette
  return tapeVariants[(props.note.id ?? 0) % tapeVariants.length]
})
const tapeStyleVariants = [
  { top: '-10px', left: '24px', transform: 'rotate(-12deg)' },
  { top: '-9px',  left: '40px', transform: 'rotate(6deg)' },
  { top: '-9px',  right: '24px', transform: 'rotate(8deg)' },
  { top: '-11px', right: '40px', transform: 'rotate(-7deg)' },
]
const tapeStyle = computed(() => tapeStyleVariants[(props.note.id ?? 0) % tapeStyleVariants.length])

const dateLabel = computed(() => {
  if (!props.note.created_at) return ''
  const d = new Date(props.note.created_at)
  if (Number.isNaN(d.getTime())) return ''
  return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' })
})
</script>

<style scoped lang="scss">
.sb-note-art {
  position: relative;
}

.na-head {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}
.na-head .r-tag {
  position: relative;
  width: 22px;
  height: 22px;
}
.na-head .r-tag.has-photo {
  background: #fdfcf7;
  overflow: hidden;
}
.na-head .r-ph {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.na-name {
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
  color: var(--sb-ink);
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.na-date {
  font-family: var(--sb-hand);
  font-size: 13px;
  color: var(--sb-ink-soft);
  white-space: nowrap;
}

.na-text {
  font-family: var(--sb-hand);
  font-size: 18px;
  line-height: 1.35;
  color: var(--sb-ink);
  margin: 0;
  word-break: break-word;
  white-space: pre-wrap;

  &.struck {
    text-decoration: line-through;
    color: var(--sb-ink-mute);
  }
}

.na-foot {
  margin-top: 10px;
  padding-top: 8px;
  border-top: 1px dashed rgba(40, 30, 20, 0.18);
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
}
.na-link {
  color: var(--sb-ink);
  text-decoration: none;
  &:hover { color: var(--sb-terracotta); }
}
.na-city { color: var(--sb-ink-mute); }

.na-x {
  position: absolute;
  top: -8px;
  right: -8px;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: var(--sb-ink);
  color: #fdfcf7;
  border: 2px solid #fdfcf7;
  font-family: var(--sb-serif);
  font-size: 14px;
  line-height: 1;
  cursor: pointer;
  z-index: 2;
  box-shadow: 0 2px 4px rgba(40, 30, 20, 0.2);
  &:hover { background: var(--sb-terracotta); }
}
</style>
