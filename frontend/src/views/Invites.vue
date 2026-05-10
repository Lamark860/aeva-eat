<template>
  <div class="sb-paper sb-grain sb-screen invites">
    <header class="invites-head">
      <h1 class="title">Приглашения</h1>
      <p class="sub">пускать в круг — по одной ссылке за раз</p>
    </header>

    <section class="generator">
      <h2 class="settings-h">Создать приглашение</h2>
      <p class="hand-line">ссылка позволит новому человеку зарегистрироваться</p>

      <button class="btn-apply" type="button" @click="createInvite" :disabled="creating">
        <span v-if="creating" class="spinner-border spinner-border-sm me-1"></span>
        сгенерировать ссылку
      </button>

      <div v-if="newInviteUrl" class="generated">
        <input type="text" class="form-control paper-control" :value="newInviteUrl" readonly />
        <button class="copy-btn" type="button" @click="copyLink">
          {{ copied ? '✓ скопировано' : 'копировать' }}
        </button>
      </div>
    </section>

    <section class="invites-list">
      <h2 class="settings-h">Мои приглашения</h2>

      <div v-if="loading" class="sb-empty">…</div>
      <div v-else-if="!invites.length" class="sb-empty">пока нет приглашений</div>

      <div v-else class="tickets">
        <div
          v-for="(inv, i) in invites"
          :key="inv.id"
          class="invite-ticket"
          :class="ticketTilt(i)"
        >
          <div class="stub code-stub">
            <span class="lbl">код</span>
            <span class="code">{{ inv.code }}</span>
          </div>
          <div class="stub status-stub">
            <Stamp :kind="inv.used_by ? 'ink' : 'gem'">
              {{ inv.used_by ? 'использован' : 'активен' }}
            </Stamp>
            <span class="date">{{ formatDate(inv.created_at) }}</span>
          </div>
          <button
            v-if="!inv.used_by"
            class="remove"
            type="button"
            :aria-label="`Удалить инвайт ${inv.code}`"
            @click="deleteInvite(inv.id)"
          >
×
</button>
        </div>
      </div>
    </section>

    <div class="back-row">
      <router-link to="/profile" class="cta-link">← в профиль</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Stamp from '../components/scrapbook/Stamp.vue'
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
  if (!confirm('Удалить инвайт?')) return
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

const tilts = ['sb-t-l1', 'sb-t-r1', 'sb-t-l2', 'sb-t-r2']
function ticketTilt(i) { return tilts[i % tilts.length] }
</script>

<style scoped lang="scss">
.invites {
  padding-top: calc(18px + var(--aeva-safe-top, 0px));
  padding-inline: 16px;
  max-width: 640px;
  margin: 0 auto;
}

.invites-head {
  text-align: center;
  margin-bottom: 18px;

  .title {
    font-family: var(--sb-serif);
    font-style: italic;
    font-weight: 500;
    font-size: 26px;
    color: var(--sb-ink);
    margin: 0;
  }
  .sub {
    font-family: var(--sb-hand);
    font-size: 16px;
    color: var(--sb-ink-mute);
    margin: 4px 0 0;
  }
}

.settings-h {
  font-family: var(--sb-serif);
  font-style: italic;
  font-weight: 500;
  font-size: 18px;
  color: var(--sb-ink);
  margin: 0 0 6px;
}

.generator {
  background: var(--sb-paper-card);
  padding: 14px 16px;
  border-radius: 1px;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.06),
    0 4px 14px rgba(40, 30, 20, 0.07);
  margin-bottom: 20px;

  .hand-line {
    font-family: var(--sb-hand);
    font-size: 16px;
    color: var(--sb-ink-mute);
    margin: 0 0 10px;
  }
}

.btn-apply {
  background: var(--sb-terracotta);
  color: var(--sb-on-accent);
  border: none;
  border-radius: 999px;
  padding: 10px 18px;
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 15px;
  cursor: pointer;
  &:disabled { opacity: 0.5; cursor: not-allowed; }
  &:hover:not(:disabled) { background: oklch(0.55 0.14 30); }
}

.generated {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;

  .paper-control {
    flex: 1 1 240px;
    font-family: var(--sb-serif);
    font-size: 14px;
    color: var(--sb-ink);
    background: oklch(0.97 0.018 82);
    border: 1.4px solid rgba(40, 30, 20, 0.18);
    border-radius: 3px;
    padding: 8px 10px;
    box-shadow: inset 0 1px 2px rgba(40, 30, 20, 0.04);
  }
}
.copy-btn {
  background: oklch(0.93 0.04 85);
  border: none;
  border-radius: 999px;
  padding: 8px 14px;
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
  color: var(--sb-ink);
  cursor: pointer;
  box-shadow: 0 1px 1px rgba(40, 30, 20, 0.08);
  &:hover { background: oklch(0.92 0.07 25); color: var(--sb-terracotta); }
}

.invites-list { margin-bottom: 24px; }

.tickets {
  display: flex;
  flex-direction: column;
  gap: 14px;
  margin-top: 12px;
}

.invite-ticket {
  display: flex;
  align-items: stretch;
  background: oklch(0.93 0.04 85);
  font-family: var(--sb-serif);
  position: relative;
  border-radius: 2px;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.08),
    0 2px 6px rgba(40, 30, 20, 0.08);

  /* punched circles */
  &::before, &::after {
    content: '';
    position: absolute;
    width: 12px; height: 12px;
    border-radius: 50%;
    background: var(--sb-paper);
    top: 50%;
    transform: translateY(-50%);
    box-shadow: inset 0 0 0 0.5px rgba(40, 30, 20, 0.1);
  }
  &::before { left: -6px; }
  &::after  { right: -6px; }

  .stub {
    padding: 10px 14px;
    display: flex;
    flex-direction: column;
    gap: 4px;
    border-right: 1.5px dashed rgba(40, 30, 20, 0.25);

    &:last-of-type { border-right: none; }

    .lbl {
      font-size: 9px;
      letter-spacing: 0.18em;
      text-transform: uppercase;
      color: var(--sb-ink-mute);
      line-height: 1;
    }
    .code {
      font-family: var(--sb-hand);
      font-size: 22px;
      color: var(--sb-ink);
      line-height: 1;
      letter-spacing: 0.04em;
      word-break: break-all;
      overflow-wrap: anywhere;
    }
    .date {
      font-family: var(--sb-hand);
      font-size: 14px;
      color: var(--sb-ink-mute);
    }
  }

  .code-stub { flex: 1 1 auto; }
  .status-stub {
    align-items: flex-start;
    flex: 0 0 auto;
  }

  .remove {
    position: absolute;
    top: -8px;
    right: -8px;
    width: 22px;
    height: 22px;
    border-radius: 50%;
    background: oklch(0.92 0.07 25);
    color: var(--sb-terracotta);
    border: none;
    font-size: 16px;
    line-height: 1;
    cursor: pointer;
    box-shadow: 0 1px 2px rgba(40, 15, 5, 0.3);
    z-index: 4;
    &:hover { background: var(--sb-terracotta); color: var(--sb-on-accent); }
  }
}

.back-row {
  text-align: center;
  padding: 8px 0 24px;
}

.cta-link {
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
  color: var(--sb-ink-mute);
  text-decoration: none;
  &:hover { color: var(--sb-terracotta); }
}
</style>
