<template>
  <div class="video-recorder">
    <!-- Preview / Result -->
    <div class="video-circle" :class="{ recording: isRecording }">
      <video
        v-if="!videoBlob"
        ref="preview"
        autoplay
        muted
        playsinline
        class="video-preview"
      ></video>
      <video
        v-else
        ref="playback"
        :src="videoUrl"
        class="video-preview"
        playsinline
        loop
        @click="togglePlayback"
      ></video>

      <!-- Recording timer -->
      <div v-if="isRecording" class="recording-timer">
        <span class="rec-dot"></span> {{ formatTime(elapsed) }}
      </div>
    </div>

    <!-- Controls -->
    <div class="video-controls mt-2">
      <template v-if="!videoBlob">
        <button
          v-if="!isRecording && !cameraReady"
          class="btn btn-sm btn-outline-primary"
          @click="startCamera"
          :disabled="starting"
        >
          <span v-if="starting" class="spinner-border spinner-border-sm me-1"></span>
          📹 Записать видео
        </button>
        <button
          v-if="cameraReady && !isRecording"
          class="btn btn-sm btn-danger"
          @click="startRecording"
        >
          ⏺ Начать запись
        </button>
        <button
          v-if="isRecording"
          class="btn btn-sm btn-outline-danger"
          @click="stopRecording"
        >
          ⏹ Стоп ({{ formatTime(maxDuration - elapsed) }})
        </button>
      </template>
      <template v-else>
        <button class="btn btn-sm btn-outline-secondary" @click="reset">
          🔄 Переснять
        </button>
        <button class="btn btn-sm btn-primary" @click="$emit('recorded', videoBlob)" :disabled="uploading">
          <span v-if="uploading" class="spinner-border spinner-border-sm me-1"></span>
          💾 Сохранить
        </button>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, onBeforeUnmount } from 'vue'

const props = defineProps({
  maxDuration: { type: Number, default: 60 },
  uploading: { type: Boolean, default: false }
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

let stream = null
let recorder = null
let chunks = []
let timer = null

async function startCamera() {
  starting.value = true
  try {
    stream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'user', width: { ideal: 480 }, height: { ideal: 480 } },
      audio: true
    })
    cameraReady.value = true
    // Wait for the preview element to be rendered
    await new Promise(r => setTimeout(r, 50))
    if (preview.value) {
      preview.value.srcObject = stream
    }
  } catch (e) {
    console.error('Camera access denied:', e)
  } finally {
    starting.value = false
  }
}

function startRecording() {
  if (!stream) return

  chunks = []
  const mimeType = MediaRecorder.isTypeSupported('video/webm;codecs=vp9,opus')
    ? 'video/webm;codecs=vp9,opus'
    : 'video/webm'

  recorder = new MediaRecorder(stream, { mimeType })
  recorder.ondataavailable = (e) => {
    if (e.data.size > 0) chunks.push(e.data)
  }
  recorder.onstop = () => {
    videoBlob.value = new Blob(chunks, { type: 'video/webm' })
    videoUrl.value = URL.createObjectURL(videoBlob.value)
    stopCamera()
    // Auto-play the result
    setTimeout(() => {
      if (playback.value) playback.value.play()
    }, 100)
  }

  recorder.start()
  isRecording.value = true
  elapsed.value = 0
  timer = setInterval(() => {
    elapsed.value++
    if (elapsed.value >= props.maxDuration) {
      stopRecording()
    }
  }, 1000)
}

function stopRecording() {
  if (recorder && recorder.state !== 'inactive') {
    recorder.stop()
  }
  isRecording.value = false
  clearInterval(timer)
}

function stopCamera() {
  if (stream) {
    stream.getTracks().forEach(t => t.stop())
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
  if (playback.value.paused) {
    playback.value.play()
  } else {
    playback.value.pause()
  }
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

<style scoped>
.video-circle {
  width: 200px;
  height: 200px;
  border-radius: 50%;
  overflow: hidden;
  position: relative;
  background: #000;
  margin: 0 auto;
  border: 3px solid #dee2e6;
  transition: border-color 0.3s;
}
.video-circle.recording {
  border-color: #dc3545;
  animation: pulse-border 1.5s infinite;
}
.video-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
  cursor: pointer;
}
.recording-timer {
  position: absolute;
  bottom: 8px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0,0,0,0.6);
  color: #fff;
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 0.75rem;
  display: flex;
  align-items: center;
  gap: 4px;
}
.rec-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #dc3545;
  animation: blink 1s infinite;
}
.video-controls {
  text-align: center;
  display: flex;
  justify-content: center;
  gap: 8px;
}
@keyframes pulse-border {
  0%, 100% { border-color: #dc3545; }
  50% { border-color: #ff6b6b; }
}
@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}
</style>
