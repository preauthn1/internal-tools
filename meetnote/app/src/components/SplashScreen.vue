<script setup lang="ts">
import { onMounted, onUnmounted, ref } from "vue";

const slow = ref(false);
let timer: ReturnType<typeof setTimeout> | undefined;
onMounted(() => (timer = setTimeout(() => (slow.value = true), 8000)));
onUnmounted(() => clearTimeout(timer));
</script>

<template>
  <div class="splash" role="status" aria-live="polite">
    <div class="wave" aria-hidden="true">
      <span v-for="i in 7" :key="i" :style="{ animationDelay: `${(i - 1) * 0.11}s` }" />
    </div>
    <h1>會議筆記</h1>
    <p>正在載入辨識模型…</p>
    <p class="muted" :class="{ shown: slow }">更新後第一次啟動需要準備 GPU，會比平常久一點。</p>
  </div>
</template>

<style scoped>
.splash {
  position: fixed;
  inset: 0;
  z-index: 100;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  background: var(--bg);
  animation: fade-in 0.3s ease-out;
}

.wave {
  display: flex;
  align-items: center;
  gap: 6px;
  height: 56px;
  margin-bottom: 20px;
}

.wave span {
  width: 6px;
  height: 100%;
  border-radius: 3px;
  background: var(--accent);
  animation: bar 1s ease-in-out infinite;
  transform-origin: center;
}

h1 {
  font-size: 22px;
  margin: 0 0 4px;
}

p {
  margin: 0;
}

.muted {
  font-size: 13px;
  margin-top: 8px;
  opacity: 0;
  transition: opacity 0.4s;
}

.muted.shown {
  opacity: 1;
}

@keyframes bar {
  0%,
  100% {
    transform: scaleY(0.25);
    opacity: 0.6;
  }
  50% {
    transform: scaleY(1);
    opacity: 1;
  }
}

@keyframes fade-in {
  from {
    opacity: 0;
  }
}

@media (prefers-reduced-motion: reduce) {
  .wave span {
    animation: none;
    transform: scaleY(0.6);
  }
}
</style>
