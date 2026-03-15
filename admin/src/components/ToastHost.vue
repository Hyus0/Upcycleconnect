<template>
  <div class="toast-host">
    <article v-for="toast in toasts" :key="toast.id" class="toast" :class="`tone-${toast.tone}`">
      <div>
        <strong>{{ toast.title }}</strong>
        <p>{{ toast.message }}</p>
      </div>
      <button type="button" class="close" @click="dismissToast(toast.id)">x</button>
    </article>
  </div>
</template>

<script setup>
import { useToastStore } from "../store/toastStore";

const { toasts, dismissToast } = useToastStore();
</script>

<style scoped>
.toast-host {
  position: fixed;
  top: 22px;
  right: 22px;
  z-index: 1000;
  display: grid;
  gap: 12px;
  width: min(360px, calc(100vw - 32px));
}

.toast {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 14px;
  padding: 16px 18px;
  border-radius: 18px;
  box-shadow: var(--shadow-soft);
  background: rgba(255, 255, 255, 0.97);
  border: 1px solid var(--border);
}

.toast p {
  margin: 6px 0 0;
  color: var(--text-secondary);
}

.tone-green {
  border-left: 5px solid var(--brand-green);
}

.tone-amber {
  border-left: 5px solid var(--accent-amber);
}

.tone-coral {
  border-left: 5px solid var(--accent-coral);
}

.close {
  border: none;
  background: transparent;
  color: var(--text-muted);
}
</style>
