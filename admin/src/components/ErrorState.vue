<template>
  <div class="surface-card state-card">
    <div class="icon">!</div>
    <div>
      <strong>{{ title }}</strong>
      <p>{{ message }}</p>
      <button v-if="retryLabel" class="button button-danger" @click="$emit('retry')">
        {{ retryLabel }}
      </button>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    default: "Erreur de chargement"
  },
  message: {
    type: String,
    default: "Une erreur est survenue."
  },
  retryLabel: {
    type: String,
    default: ""
  }
});

defineEmits(["retry"]);
</script>

<style scoped>
.state-card {
  padding: 22px;
  display: flex;
  gap: 14px;
  align-items: flex-start;
}

.icon {
  width: 36px;
  height: 36px;
  border-radius: 12px;
  display: grid;
  place-items: center;
  font-weight: 800;
  color: white;
  background: var(--accent-coral);
}

p {
  margin: 6px 0 12px;
  color: var(--text-secondary);
}
</style>
