<template>
  <Teleport to="body">
    <transition name="sb-fade">
      <div v-if="open" class="sb-sheet-backdrop" @click.self="close">
        <transition name="sb-slide" appear>
          <div v-if="open" class="sb-sheet" role="dialog" aria-modal="true" aria-label="что добавить">
            <div class="grip" aria-hidden="true"></div>
            <div class="title">что прикнопить?</div>

            <button class="item" type="button" @click="pick('visit')">
              <span class="glyph">✦</span>
              <span class="label">новый визит</span>
              <span class="hint">фото · оценки · видео</span>
            </button>

            <button class="item disabled" type="button" disabled>
              <span class="glyph">✎</span>
              <span class="label">записка от руки</span>
              <span class="hint">скоро</span>
            </button>

            <button class="item disabled" type="button" disabled>
              <span class="glyph">✿</span>
              <span class="label">в&nbsp;wishlist</span>
              <span class="hint">скоро</span>
            </button>

            <button class="cancel" type="button" @click="close">отмена</button>
          </div>
        </transition>
      </div>
    </transition>
  </Teleport>
</template>

<script setup>
import { watch } from 'vue'

const props = defineProps({
  open: { type: Boolean, default: false },
})
const emit = defineEmits(['update:open', 'pick'])

function close() { emit('update:open', false) }
function pick(kind) {
  emit('pick', kind)
  close()
}

// lock body scroll while sheet is open
watch(() => props.open, (v) => {
  if (typeof document === 'undefined') return
  document.body.style.overflow = v ? 'hidden' : ''
})
</script>

<style scoped>
.sb-fade-enter-active,
.sb-fade-leave-active { transition: opacity 0.2s ease; }
.sb-fade-enter-from,
.sb-fade-leave-to     { opacity: 0; }

.sb-slide-enter-active,
.sb-slide-leave-active { transition: transform 0.25s cubic-bezier(0.2, 0.7, 0.3, 1); }
.sb-slide-enter-from,
.sb-slide-leave-to     { transform: translateY(100%); }
</style>
