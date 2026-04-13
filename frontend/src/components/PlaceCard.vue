<template>
  <div class="card h-100 shadow-sm">
    <img
      v-if="place.image_url"
      :src="place.image_url"
      :alt="place.name"
      class="card-img-top"
      style="height: 180px; object-fit: cover"
    />
    <div class="card-body">
      <h5 class="card-title">{{ place.name }}</h5>
      <p class="text-muted mb-1" v-if="place.cuisine_type">
        <small>🍽 {{ place.cuisine_type }}</small>
      </p>
      <p class="text-muted mb-2" v-if="place.city">
        <small>📍 {{ place.city }}<span v-if="place.address">, {{ place.address }}</span></small>
      </p>

      <div class="d-flex gap-2 mb-2" v-if="place.review_count > 0">
        <span class="badge bg-warning text-dark" title="Кухня">🍴 {{ avgRound(place.avg_food) }}</span>
        <span class="badge bg-info text-dark" title="Сервис">🤝 {{ avgRound(place.avg_service) }}</span>
        <span class="badge bg-success" title="Вайб">✨ {{ avgRound(place.avg_vibe) }}</span>
        <span class="badge bg-secondary">{{ place.review_count }} отз.</span>
      </div>
      <div v-else>
        <span class="badge bg-light text-muted">Нет отзывов</span>
      </div>

      <router-link :to="`/places/${place.id}`" class="btn btn-outline-primary btn-sm mt-2">
        Подробнее →
      </router-link>
    </div>
  </div>
</template>

<script setup>
defineProps({
  place: { type: Object, required: true }
})

function avgRound(val) {
  return val ? Number(val).toFixed(1) : '–'
}
</script>
