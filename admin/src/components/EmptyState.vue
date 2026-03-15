<template>
  <div class="surface-card state-card">
    <strong>{{ title }}</strong>
    <p>{{ message }}</p>
    <slot />
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    default: "Aucune donnee"
  },
  message: {
    type: String,
    default: "Aucun element ne correspond aux filtres actuels."
  }
});
</script>

<style scoped>
.state-card {
  padding: 22px;
}

p {
  margin: 8px 0 0;
  color: var(--text-secondary);
}
</style>
