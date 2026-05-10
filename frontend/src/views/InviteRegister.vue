<template>
  <div class="sb-paper sb-grain sb-screen auth">
    <header class="auth-brand">
      <div class="sb-wordmark">aeva<span class="dot"></span>eat</div>
      <div class="sub">камерный дневник еды</div>
    </header>

    <div class="auth-card">
      <h1 class="title">Регистрация</h1>

      <div v-if="checking" class="auth-stage">
        <div class="spinner-border spinner-border-sm" role="status"></div>
        <span class="hand-line">проверяем приглашение…</span>
      </div>

      <div v-else-if="invalidCode" class="auth-stage">
        <div class="bad-stamp">
          <Stamp kind="ink">недействительно</Stamp>
        </div>
        <p class="hand-line">{{ invalidMessage }}</p>
        <router-link to="/login" class="cta-link">← к входу</router-link>
      </div>

      <template v-else>
        <div class="invite-banner">
          <span class="hand-line">вас пригласил(а)</span>
          <Stamp kind="moss" class="creator">{{ creatorName }}</Stamp>
        </div>

        <div v-if="error" class="auth-error">{{ error }}</div>

        <form @submit.prevent="handleRegister">
          <div class="field">
            <label for="r-username" class="lbl">логин *</label>
            <input
              id="r-username"
              v-model="username"
              type="text"
              class="form-control paper-control"
              placeholder="username"
              required
              autocomplete="username"
            />
          </div>

          <div class="field">
            <label for="r-display" class="lbl">имя (необязательно)</label>
            <input
              id="r-display"
              v-model="displayName"
              type="text"
              class="form-control paper-control"
              placeholder="как вас называть"
            />
          </div>

          <div class="field">
            <label for="r-password" class="lbl">пароль *</label>
            <input
              id="r-password"
              v-model="password"
              type="password"
              class="form-control paper-control"
              placeholder="минимум 6 символов"
              minlength="6"
              required
              autocomplete="new-password"
            />
          </div>

          <button type="submit" class="btn-apply" :disabled="loading">
            <span v-if="loading" class="spinner-border spinner-border-sm me-1"></span>
            создать аккаунт
          </button>
        </form>

        <p class="hint">
          уже есть аккаунт?
          <router-link to="/login" class="cta-link">войти</router-link>
        </p>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import Stamp from '../components/scrapbook/Stamp.vue'
import http from '../api/http'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const inviteCode = route.params.code

const checking = ref(true)
const invalidCode = ref(false)
const invalidMessage = ref('')
const creatorName = ref('')

const username = ref('')
const displayName = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

onMounted(async () => {
  try {
    const { data } = await http.get(`/invites/validate/${inviteCode}`)
    if (data.valid) {
      creatorName.value = data.creator_name
    } else {
      invalidCode.value = true
      invalidMessage.value = 'этот код не работает'
    }
  } catch {
    // Backend errors here are typically English ('invalid invite code',
    // 'expired', 'already used'). Show a friendly Russian line instead.
    invalidCode.value = true
    invalidMessage.value = 'этот код не работает'
  } finally {
    checking.value = false
  }
})

async function handleRegister() {
  error.value = ''
  loading.value = true
  try {
    await auth.register(username.value, displayName.value, password.value, inviteCode)
    router.push('/')
  } catch (e) {
    error.value = e.response?.data?.error || 'Ошибка регистрации'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.auth {
  padding: calc(40px + var(--aeva-safe-top, 0px)) 16px 0;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.auth-brand {
  text-align: center;
  margin-bottom: 28px;

  .sb-wordmark { font-size: 30px; }
  .sub {
    font-family: var(--sb-hand);
    font-size: 18px;
    color: var(--sb-ink-mute);
    margin-top: 4px;
  }
}

.auth-card {
  background: var(--sb-paper-card);
  width: 100%;
  max-width: 360px;
  padding: 22px 22px 18px;
  position: relative;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.08),
    0 6px 14px rgba(40, 30, 20, 0.12),
    0 0 0 0.5px rgba(40, 30, 20, 0.06);
  border-radius: 1px;
  transform: rotate(0.6deg);
}

.title {
  font-family: var(--sb-serif);
  font-style: italic;
  font-weight: 500;
  font-size: 24px;
  color: var(--sb-ink);
  text-align: center;
  margin: 0 0 14px;
}

.auth-stage {
  text-align: center;
  padding: 18px 4px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: center;
}

.bad-stamp { display: inline-block; transform: rotate(-2deg); }

.invite-banner {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-bottom: 16px;
  padding: 10px 12px;
  background: oklch(0.94 0.05 85 / 0.5);
  border: 1px dashed rgba(40, 30, 20, 0.25);
  border-radius: 4px;
}
.invite-banner .creator { transform: rotate(-1.5deg); }

.hand-line {
  font-family: var(--sb-hand);
  font-size: 16px;
  color: var(--sb-ink-soft);
  line-height: 1.2;
}

.auth-error {
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
  background: oklch(0.92 0.07 25);
  color: var(--sb-terracotta);
  border: 1px dashed oklch(0.55 0.14 30 / 0.5);
  border-radius: 4px;
  padding: 8px 12px;
  margin-bottom: 12px;
}

.field { margin-bottom: 12px; }

.lbl {
  display: block;
  font-family: var(--sb-hand);
  font-size: 16px;
  color: var(--sb-ink-mute);
  margin-bottom: 4px;
}

.paper-control {
  font-family: var(--sb-serif);
  font-size: 16px;
  color: var(--sb-ink);
  background: oklch(0.97 0.018 82);
  border: 1.4px solid rgba(40, 30, 20, 0.18);
  border-radius: 3px;
  padding: 10px 12px;
  width: 100%;
  box-shadow: inset 0 1px 2px rgba(40, 30, 20, 0.04);

  &:focus {
    outline: none;
    border-color: var(--sb-terracotta);
    box-shadow: 0 0 0 2px oklch(0.55 0.14 30 / 0.15);
    background: var(--sb-paper-card);
  }
  &::placeholder {
    color: var(--sb-ink-mute);
    font-style: italic;
  }
}

.btn-apply {
  background: var(--sb-terracotta);
  color: var(--sb-on-accent);
  border: none;
  border-radius: 999px;
  padding: 12px 22px;
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 16px;
  cursor: pointer;
  width: 100%;
  margin-top: 6px;
  &:disabled { opacity: 0.5; cursor: not-allowed; }
  &:hover:not(:disabled) { background: oklch(0.55 0.14 30); }
}

.hint {
  font-family: var(--sb-hand);
  font-size: 16px;
  color: var(--sb-ink-mute);
  text-align: center;
  margin: 14px 0 0;

  .cta-link { color: var(--sb-terracotta); }
}

.cta-link {
  font-family: var(--sb-serif);
  font-style: italic;
  text-decoration: none;
  color: var(--sb-ink);
  &:hover { color: var(--sb-terracotta); }
}
</style>
