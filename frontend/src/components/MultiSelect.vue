<template>
  <div class="sb-multi" ref="wrapper">
    <div
      class="trigger"
      :class="{ open }"
      role="button"
      tabindex="0"
      @click="toggleOpen"
      @keydown.enter.prevent="toggleOpen"
      @keydown.space.prevent="toggleOpen"
    >
      <span
        v-for="item in selectedItems"
        :key="item.id"
        class="chip"
      >
        {{ item.name }}
        <span class="chip-x" role="button" aria-label="Убрать" @click.stop="toggle(item.id)">×</span>
      </span>
      <span v-if="selectedItems.length === 0" class="ph">{{ placeholder }}</span>
      <span class="caret" aria-hidden="true">▾</span>
    </div>

    <div v-if="open" class="dropdown" role="listbox">
      <div
        v-for="option in options"
        :key="option.id"
        class="opt"
        :class="{ active: modelValue.includes(option.id) }"
        role="option"
        :aria-selected="modelValue.includes(option.id)"
        @click="toggle(option.id)"
      >
        <span class="check">{{ modelValue.includes(option.id) ? '✓' : '' }}</span>
        {{ option.name }}
      </div>
      <div v-if="options.length === 0" class="empty">пусто</div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  modelValue:  { type: Array,  default: () => [] },
  options:     { type: Array,  default: () => [] },
  placeholder: { type: String, default: 'выберите…' },
})

const emit = defineEmits(['update:modelValue'])

const open = ref(false)
const wrapper = ref(null)

const selectedItems = computed(() =>
  props.options.filter((o) => props.modelValue.includes(o.id)),
)

function toggleOpen() { open.value = !open.value }

function toggle(id) {
  const current = [...props.modelValue]
  const idx = current.indexOf(id)
  if (idx >= 0) current.splice(idx, 1)
  else current.push(id)
  emit('update:modelValue', current)
}

function onClickOutside(e) {
  if (wrapper.value && !wrapper.value.contains(e.target)) {
    open.value = false
  }
}

onMounted(() => document.addEventListener('click', onClickOutside))
onUnmounted(() => document.removeEventListener('click', onClickOutside))
</script>

<style scoped lang="scss">
.sb-multi {
  position: relative;
  font-family: var(--sb-serif);
}

.trigger {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 6px;
  min-height: 40px;
  padding: 6px 30px 6px 10px;
  background: oklch(0.97 0.018 82);
  border: 1.4px solid rgba(40, 30, 20, 0.18);
  border-radius: 3px;
  font-family: var(--sb-serif);
  font-size: 15px;
  color: var(--sb-ink);
  cursor: pointer;
  position: relative;
  box-shadow: inset 0 1px 2px rgba(40, 30, 20, 0.04);
  outline: none;

  &:focus-visible,
  &.open {
    border-color: var(--sb-terracotta);
    box-shadow: 0 0 0 2px oklch(0.55 0.14 30 / 0.15);
    background: var(--sb-paper-card);
  }
}

.ph {
  font-style: italic;
  color: var(--sb-ink-mute);
  flex: 1;
  line-height: 1.4;
}

.caret {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--sb-ink-mute);
  font-size: 12px;
  pointer-events: none;
  line-height: 1;
}
.trigger.open .caret { color: var(--sb-terracotta); }

.chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: transparent;
  border: 1.4px solid var(--sb-terracotta);
  color: var(--sb-terracotta);
  border-radius: 999px;
  padding: 2px 6px 2px 10px;
  font-family: var(--sb-serif);
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  line-height: 1.4;
  box-shadow: inset 0 0 0 0.5px rgba(140, 60, 30, 0.2);
  white-space: nowrap;
}
.chip-x {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  font-size: 14px;
  font-weight: 400;
  text-transform: none;
  line-height: 1;
  cursor: pointer;
  &:hover {
    background: var(--sb-terracotta);
    color: var(--sb-on-accent);
  }
}

.dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  z-index: 100;
  background: var(--sb-paper-card);
  border: 1px solid rgba(40, 30, 20, 0.18);
  border-radius: 3px;
  max-height: 220px;
  overflow-y: auto;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.06),
    0 6px 18px rgba(40, 30, 20, 0.18);
}

.opt {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  cursor: pointer;
  font-family: var(--sb-serif);
  font-size: 15px;
  color: var(--sb-ink);
  border-bottom: 1px dashed rgba(40, 30, 20, 0.1);
  transition: background 0.12s;

  &:last-child { border-bottom: none; }
  &:hover { background: oklch(0.94 0.05 85 / 0.5); }
  &.active {
    background: oklch(0.94 0.05 85 / 0.7);
    font-weight: 500;
  }
}

.check {
  width: 16px;
  text-align: center;
  font-family: var(--sb-hand);
  font-size: 18px;
  color: var(--sb-terracotta);
  line-height: 1;
}

.empty {
  padding: 10px 12px;
  font-family: var(--sb-hand);
  font-size: 16px;
  color: var(--sb-ink-mute);
  text-align: center;
}
</style>
