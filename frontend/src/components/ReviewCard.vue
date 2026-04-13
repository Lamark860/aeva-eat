<template>
  <div class="card mb-3">
    <div class="card-body">
      <div class="d-flex justify-content-between align-items-start">
        <div>
          <div class="d-flex gap-2 mb-1">
            <span class="badge bg-warning text-dark" title="Кухня">🍴 {{ review.food_rating }}</span>
            <span class="badge bg-info text-dark" title="Сервис">🤝 {{ review.service_rating }}</span>
            <span class="badge bg-success" title="Вайб">✨ {{ review.vibe_rating }}</span>
            <span v-if="review.is_gem" class="badge bg-primary" title="Жемчужина">💎</span>
          </div>
          <div class="text-muted small mb-1">
            <span v-for="(author, i) in review.authors" :key="author.id">
              <span v-if="i > 0">, </span>{{ author.username }}
            </span>
            <span v-if="review.visited_at"> · {{ review.visited_at }}</span>
          </div>
        </div>
        <div v-if="canEdit" class="d-flex gap-1">
          <button class="btn btn-outline-primary btn-sm" @click="$emit('edit', review)">✏️</button>
          <button class="btn btn-outline-danger btn-sm" @click="$emit('delete', review.id)">🗑</button>
        </div>
      </div>
      <p v-if="review.comment" class="mb-0 mt-2">{{ review.comment }}</p>
    </div>
  </div>
</template>

<script setup>
defineProps({
  review: { type: Object, required: true },
  canEdit: { type: Boolean, default: false }
})

defineEmits(['edit', 'delete'])
</script>
