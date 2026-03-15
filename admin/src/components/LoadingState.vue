<template>
  <div class="surface-card state-card">
    <div class="dot"></div>
    <div>
      <strong>{{ title }}</strong>
      <p>{{ message }}</p>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    default: "Chargement"
  },
  message: {
    type: String,
    default: "Les donnees sont en cours de recuperation."
  }
});
</script>

<style scoped>
.state-card {
  padding: 22px;
  display: flex;
  gap: 14px;
  align-items: center;
}

.dot {
  width: 16px;
  height: 16px;
  border-radius: 999px;
  background: var(--brand-green);
  box-shadow: 0 0 0 10px rgba(45, 122, 79, 0.12);
  animation: pulse 1.1s infinite ease-in-out;
}

@keyframes pulse {
  50% {
    transform: scale(1.12);
    opacity: 0.72;
  }
}

p {
  margin: 6px 0 0;
  color: var(--text-secondary);
}
</style>
