<template>
  <div class="container">
    <h2 class="mb-4">Приглашения</h2>

    <div class="card shadow-sm mb-4">
      <div class="card-body">
        <h5 class="card-title">Создать приглашение</h5>
        <p class="text-muted small">Ссылка-приглашение позволит новому пользователю зарегистрироваться</p>

        <button class="btn btn-primary" @click="createInvite" :disabled="creating">
          <span v-if="creating" class="spinner-border spinner-border-sm me-1"></span>
          Сгенерировать ссылку
        </button>

        <div v-if="newInviteUrl" class="mt-3">
          <div class="input-group">
            <input type="text" class="form-control" :value="newInviteUrl" readonly />
            <button class="btn btn-outline-secondary" @click="copyLink">
              {{ copied ? '✓ Скопировано' : '📋 Копировать' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="card shadow-sm">
      <div class="card-body">
        <h5 class="card-title">Мои приглашения</h5>

        <div v-if="loading" class="text-center py-3">
          <div class="spinner-border text-primary" role="status"></div>
        </div>

        <div v-else-if="!invites.length" class="text-muted py-3">
          Пока нет приглашений
        </div>

        <div v-else class="table-responsive">
          <table class="table table-sm align-middle">
            <thead>
              <tr>
                <th>Код</th>
                <th>Статус</th>
                <th>Создан</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="inv in invites" :key="inv.id">
                <td>
                  <code class="small">{{ inv.code }}</code>
                </td>
                <td>
                  <span v-if="inv.used_by" class="badge bg-secondary">
                    Использован
                  </span>
                  <span v-else class="badge bg-success">Активен</span>
                </td>
                <td class="small text-muted">
                  {{ formatDate(inv.created_at) }}
                </td>
                <td>
                  <button
                    v-if="!inv.used_by"
                    class="btn btn-sm btn-outline-danger"
                    @click="deleteInvite(inv.id)"
                  >
                    ✕
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import http from '../api/http'

const invites = ref([])
const loading = ref(true)
const creating = ref(false)
const newInviteUrl = ref('')
const copied = ref(false)

onMounted(() => fetchInvites())

async function fetchInvites() {
  loading.value = true
  try {
    const { data } = await http.get('/invites')
    invites.value = data || []
  } finally {
    loading.value = false
  }
}

async function createInvite() {
  creating.value = true
  copied.value = false
  try {
    const { data } = await http.post('/invites')
    newInviteUrl.value = `${window.location.origin}/invite/${data.code}`
    await fetchInvites()
  } finally {
    creating.value = false
  }
}

async function deleteInvite(id) {
  await http.delete(`/invites/${id}`)
  await fetchInvites()
}

function copyLink() {
  navigator.clipboard.writeText(newInviteUrl.value)
  copied.value = true
  setTimeout(() => { copied.value = false }, 2000)
}

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' })
}
</script>
