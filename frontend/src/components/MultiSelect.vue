<template>
  <div class="multi-select" ref="wrapper">
    <div
      class="multi-select-trigger form-control form-control-sm d-flex align-items-center flex-wrap gap-1"
      @click="open = !open"
    >
      <span
        v-for="item in selectedItems"
        :key="item.id"
        class="multi-tag"
      >
        {{ item.name }}
        <span class="multi-tag-remove" @click.stop="toggle(item.id)">&times;</span>
      </span>
      <span v-if="selectedItems.length === 0" class="text-muted">{{ placeholder }}</span>
    </div>
    <div v-if="open" class="multi-select-dropdown">
      <div
        v-for="option in options"
        :key="option.id"
        class="multi-select-option"
        :class="{ active: modelValue.includes(option.id) }"
        @click="toggle(option.id)"
      >
        <span class="multi-check">{{ modelValue.includes(option.id) ? '✓' : '' }}</span>
        {{ option.name }}
      </div>
      <div v-if="options.length === 0" class="text-muted small p-2">Нет вариантов</div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  modelValue: { type: Array, default: () => [] },
  options: { type: Array, default: () => [] },
  placeholder: { type: String, default: 'Выберите...' }
})

const emit = defineEmits(['update:modelValue'])

const open = ref(false)
const wrapper = ref(null)

const selectedItems = computed(() =>
  props.options.filter(o => props.modelValue.includes(o.id))
)

function toggle(id) {
  const current = [...props.modelValue]
  const idx = current.indexOf(id)
  if (idx >= 0) {
    current.splice(idx, 1)
  } else {
    current.push(id)
  }
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

<style scoped>
.multi-select {
  position: relative;
}

.multi-select-trigger {
  cursor: pointer;
  min-height: 31px;
  padding: 0.2rem 0.5rem;
}

.multi-tag {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  background: var(--bs-primary);
  color: #fff;
  border-radius: 0.35rem;
  padding: 0.1rem 0.45rem;
  font-size: 0.75rem;
  font-weight: 500;
  line-height: 1.4;
}

.multi-tag-remove {
  cursor: pointer;
  font-size: 0.9rem;
  line-height: 1;
  opacity: 0.8;
}
.multi-tag-remove:hover {
  opacity: 1;
}

.multi-select-dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  z-index: 100;
  background: #fff;
  border: 1px solid #e0ddd9;
  border-radius: 0.5rem;
  margin-top: 0.25rem;
  max-height: 200px;
  overflow-y: auto;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.multi-select-option {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.4rem 0.75rem;
  cursor: pointer;
  font-size: 0.85rem;
  transition: background 0.15s;
}
.multi-select-option:hover {
  background: #f5f3f1;
}
.multi-select-option.active {
  background: #fef3ef;
  font-weight: 500;
}

.multi-check {
  width: 1.1rem;
  text-align: center;
  color: var(--bs-primary);
  font-weight: 700;
}
</style>
