<template>
  <div class="sb-recorder">
    <!-- Round video frame in scrapbook border -->
    <div class="frame" :class="{ recording: isRecording, idle: !cameraReady && !videoBlob }">
      <video
        v-if="!videoBlob"
        ref="preview"
        autoplay
        muted
        playsinline
        class="vid"
      ></video>
      <video
        v-else
        ref="playback"
        :src="videoUrl"
        class="vid"
        playsinline
        loop
        @click="togglePlayback"
      ></video>

      <!-- Idle hint -->
      <div v-if="!cameraReady && !videoBlob && !starting" class="hint">
        <span v-if="error">{{ error }}</span>
        <span v-else>тапни «записать»</span>
      </div>

      <!-- Recording timer -->
      <div v-if="isRecording" class="timer">
        <span class="dot"></span> {{ formatTime(elapsed) }}
      </div>
    </div>

    <!-- Controls -->
    <div class="controls">
      <template v-if="!videoBlob">
        <button
          v-if="!isRecording && !cameraReady"
          class="btn-paper"
          type="button"
          @click="startCamera"
          :disabled="starting"
        >
          <span v-if="starting" class="spinner-border spinner-border-sm me-1"></span>
          записать видео
        </button>
        <button
          v-if="cameraReady && !isRecording"
          class="btn-rec"
          type="button"
          @click="startRecording"
        >
          <span class="rec-glyph"></span> начать запись
        </button>
        <button
          v-if="isRecording"
          class="btn-stop"
          type="button"
          @click="stopRecording"
        >
          ◼ стоп · {{ formatTime(maxDuration - elapsed) }}
        </button>
      </template>
      <template v-else>
        <button class="btn-paper" type="button" @click="reset">переснять</button>
        <button class="btn-apply" type="button" :disabled="uploading" @click="emit('recorded', videoBlob)">
          <span v-if="uploading" class="spinner-border spinner-border-sm me-1"></span>
          сохранить
        </button>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, onBeforeUnmount } from 'vue'

const props = defineProps({
  maxDuration: { type: Number,  default: 60 },
  uploading:   { type: Boolean, default: false },
})

const emit = defineEmits(['recorded'])

const preview = ref(null)
const playback = ref(null)

const cameraReady = ref(false)
const starting = ref(false)
const isRecording = ref(false)
const elapsed = ref(0)
const videoBlob = ref(null)
const videoUrl = ref('')
const error = ref('')

let stream = null
let recorder = null
let chunks = []
let timer = null

async function startCamera() {
  starting.value = true
  error.value = ''
  // getUserMedia is gated behind a secure context: localhost or HTTPS.
  // When opened from a LAN hostname (e.g. http://*.local:8091) the browser
  // disables it silently — surface a soft hint per DESIGN-DECISIONS R2.
  // Technical detail goes to console.error for the dev to investigate.
  if (!window.isSecureContext) {
    error.value = 'кружочки пока только с компьютера'
    // eslint-disable-next-line no-console
    console.error('[VideoRecorder] insecure context — getUserMedia blocked. Open via localhost or HTTPS.')
    starting.value = false
    return
  }
  if (!navigator.mediaDevices?.getUserMedia) {
    error.value = 'браузер не поддерживает запись'
    starting.value = false
    return
  }
  try {
    stream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'user', width: { ideal: 480 }, height: { ideal: 480 } },
      audio: true,
    })
    cameraReady.value = true
    await new Promise((r) => setTimeout(r, 50))
    if (preview.value) preview.value.srcObject = stream
  } catch (e) {
    cameraReady.value = false
    if (e?.name === 'NotAllowedError') error.value = 'нет разрешения на камеру'
    else if (e?.name === 'NotFoundError') error.value = 'камера не найдена'
    else error.value = `не удалось включить камеру: ${e?.message || e}`
    // eslint-disable-next-line no-console
    console.error('[VideoRecorder] getUserMedia failed:', e)
  } finally {
    starting.value = false
  }
}

// Pick the best mimeType the browser actually supports — iOS Safari only does mp4.
function pickMimeType() {
  const candidates = [
    'video/webm;codecs=vp9,opus',
    'video/webm;codecs=vp8,opus',
    'video/webm',
    'video/mp4;codecs=avc1.42E01E,mp4a.40.2',
    'video/mp4',
  ]
  for (const c of candidates) {
    if (typeof MediaRecorder !== 'undefined' && MediaRecorder.isTypeSupported(c)) return c
  }
  return ''
}

function startRecording() {
  if (!stream) return

  chunks = []
  const mimeType = pickMimeType()

  try {
    recorder = mimeType ? new MediaRecorder(stream, { mimeType }) : new MediaRecorder(stream)
  } catch (e) {
    isRecording.value = false
    error.value = `не удалось начать запись: ${e?.message || e}`
    // eslint-disable-next-line no-console
    console.error('[VideoRecorder] MediaRecorder failed:', e)
    return
  }

  recorder.ondataavailable = (e) => {
    if (e.data.size > 0) chunks.push(e.data)
  }
  recorder.onstop = () => {
    // Strip codec parameters so the backend (which validates plain video/webm | video/mp4) accepts it.
    const recorded = (recorder && recorder.mimeType) || mimeType || 'video/webm'
    const baseType = recorded.split(';')[0] || 'video/webm'
    videoBlob.value = new Blob(chunks, { type: baseType })
    videoUrl.value = URL.createObjectURL(videoBlob.value)
    stopCamera()
    setTimeout(() => {
      if (playback.value) playback.value.play()
    }, 100)
  }

  recorder.start()
  isRecording.value = true
  elapsed.value = 0
  timer = setInterval(() => {
    elapsed.value++
    if (elapsed.value >= props.maxDuration) stopRecording()
  }, 1000)
}

function stopRecording() {
  if (recorder && recorder.state !== 'inactive') recorder.stop()
  isRecording.value = false
  clearInterval(timer)
}

function stopCamera() {
  if (stream) {
    stream.getTracks().forEach((t) => t.stop())
    stream = null
  }
  cameraReady.value = false
}

function reset() {
  if (videoUrl.value) URL.revokeObjectURL(videoUrl.value)
  videoBlob.value = null
  videoUrl.value = ''
  elapsed.value = 0
}

function togglePlayback() {
  if (!playback.value) return
  if (playback.value.paused) playback.value.play()
  else playback.value.pause()
}

function formatTime(sec) {
  const m = Math.floor(sec / 60)
  const s = sec % 60
  return `${m}:${s.toString().padStart(2, '0')}`
}

onBeforeUnmount(() => {
  stopRecording()
  stopCamera()
  if (videoUrl.value) URL.revokeObjectURL(videoUrl.value)
})
</script>

<style scoped lang="scss">
.sb-recorder {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 12px;
}

.frame {
  width: 200px;
  height: 200px;
  border-radius: 50%;
  overflow: hidden;
  position: relative;
  background: #000;
  box-shadow:
    0 0 0 3px var(--sb-paper-card),
    0 2px 4px rgba(40, 30, 20, 0.18),
    0 8px 18px rgba(40, 30, 20, 0.18);
  align-self: center;

  &.idle {
    background: oklch(0.92 0.022 82);
  }
  &.recording {
    box-shadow:
      0 0 0 3px var(--sb-paper-card),
      0 0 0 5px var(--sb-terracotta),
      0 2px 4px rgba(40, 30, 20, 0.18),
      0 8px 18px rgba(40, 30, 20, 0.18);
    animation: sb-rec-ring 1.6s ease-in-out infinite;
  }
}
.vid {
  width: 100%;
  height: 100%;
  object-fit: cover;
  cursor: pointer;
}
.hint {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--sb-hand);
  font-size: 15px;
  line-height: 1.25;
  color: var(--sb-ink-mute);
  text-align: center;
  padding: 0 18px;
}

.timer {
  position: absolute;
  bottom: 8px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(40, 30, 20, 0.65);
  color: var(--sb-on-accent);
  padding: 3px 10px;
  border-radius: 999px;
  font-family: var(--sb-serif);
  font-size: 12px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--sb-terracotta);
  animation: sb-blink 1s infinite;
}

.controls {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-self: center;
}

.btn-paper,
.btn-rec,
.btn-stop,
.btn-apply {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border: none;
  border-radius: 999px;
  font-family: var(--sb-serif);
  font-style: italic;
  font-size: 14px;
  cursor: pointer;
  min-height: 36px;
  &:disabled { opacity: 0.5; cursor: not-allowed; }
}
.btn-paper {
  background: oklch(0.93 0.04 85);
  color: var(--sb-ink);
  box-shadow: 0 1px 1px rgba(40, 30, 20, 0.08), 0 2px 6px rgba(40, 30, 20, 0.10);
  &:hover:not(:disabled) { background: oklch(0.92 0.07 25); color: var(--sb-terracotta); }
}
.btn-rec {
  background: var(--sb-terracotta);
  color: var(--sb-on-accent);
  box-shadow: 0 1px 1px rgba(140, 60, 30, 0.2), 0 3px 8px rgba(140, 60, 30, 0.25);
  &:hover:not(:disabled) { background: oklch(0.55 0.14 30); }
}
.btn-rec .rec-glyph {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #fff;
  display: inline-block;
}
.btn-stop {
  background: transparent;
  color: var(--sb-terracotta);
  border: 1.4px solid var(--sb-terracotta);
  box-shadow: inset 0 0 0 0.5px rgba(140, 60, 30, 0.2);
  &:hover { background: oklch(0.92 0.07 25); }
}
.btn-apply {
  background: var(--sb-terracotta);
  color: var(--sb-on-accent);
  box-shadow: 0 1px 1px rgba(140, 60, 30, 0.2), 0 3px 8px rgba(140, 60, 30, 0.25);
  &:hover:not(:disabled) { background: oklch(0.55 0.14 30); }
}

@keyframes sb-rec-ring {
  0%, 100% { box-shadow:
    0 0 0 3px var(--sb-paper-card),
    0 0 0 5px var(--sb-terracotta),
    0 2px 4px rgba(40, 30, 20, 0.18),
    0 8px 18px rgba(40, 30, 20, 0.18); }
  50% { box-shadow:
    0 0 0 3px var(--sb-paper-card),
    0 0 0 7px oklch(0.7 0.16 25 / 0.6),
    0 2px 4px rgba(40, 30, 20, 0.18),
    0 8px 18px rgba(40, 30, 20, 0.18); }
}
@keyframes sb-blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}
</style>
