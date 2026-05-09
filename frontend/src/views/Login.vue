<template>
  <div class="sb-paper sb-grain sb-screen auth">
    <header class="auth-brand">
      <div class="sb-wordmark">aeva<span class="dot"></span>eat</div>
      <div class="sub">камерный дневник еды</div>
    </header>

    <div class="auth-card">
      <h1 class="title">Войти</h1>

      <div v-if="error" class="auth-error">{{ error }}</div>

      <form @submit.prevent="handleLogin">
        <div class="field">
          <label for="username" class="lbl">логин</label>
          <input
            id="username"
            v-model="username"
            type="text"
            class="form-control paper-control"
            placeholder="username"
            required
            autocomplete="username"
          />
        </div>

        <div class="field">
          <label for="password" class="lbl">пароль</label>
          <input
            id="password"
            v-model="password"
            type="password"
            class="form-control paper-control"
            placeholder="••••••"
            required
            autocomplete="current-password"
          />
        </div>

        <button type="submit" class="btn-apply" :disabled="loading">
          <span v-if="loading" class="spinner-border spinner-border-sm me-1"></span>
          войти
        </button>
      </form>

      <p class="hint">доступ только по приглашению</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()

const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(username.value, password.value)
    router.push('/')
  } catch (e) {
    error.value = e.response?.data?.error || 'Ошибка входа'
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
  margin-bottom: 32px;

  .sb-wordmark { font-size: 30px; }
  .sub {
    font-family: var(--sb-hand);
    font-size: 18px;
    color: var(--sb-ink-mute);
    margin-top: 4px;
  }
}

.auth-card {
  background: #fdfcf7;
  width: 100%;
  max-width: 360px;
  padding: 22px 22px 18px;
  position: relative;
  box-shadow:
    0 1px 1px rgba(40, 30, 20, 0.08),
    0 6px 14px rgba(40, 30, 20, 0.12),
    0 0 0 0.5px rgba(40, 30, 20, 0.06);
  border-radius: 1px;
  transform: rotate(-0.6deg);
}

.title {
  font-family: var(--sb-serif);
  font-style: italic;
  font-weight: 500;
  font-size: 26px;
  color: var(--sb-ink);
  text-align: center;
  margin: 0 0 18px;
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

.field { margin-bottom: 14px; }

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
    background: #fdfcf7;
  }
  &::placeholder {
    color: var(--sb-ink-mute);
    font-style: italic;
  }
}

.btn-apply {
  background: var(--sb-terracotta);
  color: #fff;
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
}
</style>
