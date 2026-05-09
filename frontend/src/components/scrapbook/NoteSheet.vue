<template>
  <Teleport to="body">
    <transition name="sb-fade">
      <div v-if="open" class="sb-sheet-backdrop" @click.self="close">
        <transition name="sb-slide" appear>
          <div v-if="open" class="sb-sheet note-sheet" role="dialog" aria-modal="true" aria-label="новая записка">
            <div class="grip" aria-hidden="true"></div>
            <header class="ns-head">
              <h3 class="ns-title">записка от руки</h3>
              <span class="ns-sub">пиши коротко и от души</span>
            </header>

            <form class="ns-form" @submit.prevent="submit">
              <textarea
                v-model="form.text"
                class="paper-control ns-text"
                rows="4"
                maxlength="2000"
                placeholder="что прикнопить?"
                required
              ></textarea>

              <input
                v-model="form.city"
                type="text"
                class="paper-control"
                placeholder="город (опционально)"
              />

              <div class="ns-cta">
                <button type="submit" class="btn-apply" :disabled="!form.text.trim() || loading">
                  {{ loading ? '…' : 'прикнопить' }}
                </button>
                <button type="button" class="cancel-link" @click="close">отмена</button>
              </div>
            </form>
          </div>
        </transition>
      </div>
    </transition>
  </Teleport>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  open: { type: Boolean, default: false },
})
const emit = defineEmits(['update:open', 'submit'])

const form = ref({ text: '', city: '' })
const loading = ref(false)

watch(() => props.open, (v) => {
  if (typeof document !== 'undefined') {
    document.body.style.overflow = v ? 'hidden' : ''
  }
  if (v) {
    form.value = { text: '', city: '' }
    loading.value = false
  }
})

function close() {
  emit('update:open', false)
}

async function submit() {
  loading.value = true
  emit('submit', {
    text: form.value.text.trim(),
    city: form.value.city.trim() || undefined,
  })
  // родитель закроет через update:open после await
}
</script>

<style scoped lang="scss">
.note-sheet {
  /* Локальный фон/паддинги наследуются от .sb-sheet (определён глобально). */
  padding-bottom: calc(20px + var(--aeva-safe-bottom, 0px));
}

.ns-head {
  padding: 4px 18px 12px;
  text-align: center;
}
.ns-title {
  font-family: var(--sb-serif);
  font-style: italic;
  font-weight: 500;
  font-size: 22px;
  margin: 0;
  color: var(--sb-ink);
}
.ns-sub {
  font-family: var(--sb-hand);
  font-size: 16px;
  color: var(--sb-ink-mute);
}

.ns-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 0 18px 12px;
}

.paper-control {
  width: 100%;
  font-family: var(--sb-serif);
  font-size: 16px;
  color: var(--sb-ink);
  background: oklch(0.97 0.018 82);
  border: 1.4px solid rgba(40, 30, 20, 0.18);
  border-radius: 3px;
  padding: 10px 12px;
  box-shadow: inset 0 1px 2px rgba(40, 30, 20, 0.04);

  &:focus {
    border-color: var(--sb-terracotta);
    outline: none;
    box-shadow: 0 0 0 2px oklch(0.55 0.14 30 / 0.15);
    background: #fdfcf7;
  }
  &::placeholder {
    color: var(--sb-ink-mute);
    font-style: italic;
  }
}

textarea.ns-text {
  font-family: var(--sb-hand);
  font-size: 19px;
  line-height: 1.4;
  resize: vertical;
  min-height: 120px;
}

.ns-cta {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-top: 4px;
}
.btn-apply {
  background: var(--sb-terracotta);
  color: #fff;
  border: none;
  border-radius: 999px;
  padding: 12px 22px;
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 15px;
  cursor: pointer;
  flex: 1;
  &:disabled { opacity: 0.5; cursor: not-allowed; }
  &:hover:not(:disabled) { background: oklch(0.55 0.14 30); color: #fff; }
}
.cancel-link {
  background: transparent;
  border: none;
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 15px;
  color: var(--sb-ink-mute);
  cursor: pointer;
  padding: 12px 8px;
  &:hover { color: var(--sb-ink); }
}
</style>
