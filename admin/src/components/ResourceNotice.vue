<template>
  <article class="surface-card notice-card" :class="`tone-${tone}`">
    <div class="notice-head">
      <h3>{{ title }}</h3>
      <StatusBadge :label="badge" :tone="tone" />
    </div>
    <p>{{ message }}</p>
    <ul>
      <li v-for="item in items" :key="item">{{ item }}</li>
    </ul>
  </article>
</template>

<script setup>
import StatusBadge from "./StatusBadge.vue";

defineProps({
  title: {
    type: String,
    required: true
  },
  message: {
    type: String,
    required: true
  },
  items: {
    type: Array,
    default: () => []
  },
  tone: {
    type: String,
    default: "amber"
  },
  badge: {
    type: String,
    default: "a traiter"
  }
});
</script>

<style scoped>
.notice-card {
  padding: 20px;
  border-left: 6px solid transparent;
}

.tone-amber {
  border-left-color: var(--accent-amber);
}

.tone-coral {
  border-left-color: var(--accent-coral);
}

.tone-green {
  border-left-color: var(--brand-green);
}

.notice-head {
  display: flex;
  justify-content: space-between;
  gap: 14px;
  align-items: center;
}

h3 {
  margin: 0;
  font-family: "Syne", sans-serif;
}

p {
  margin: 8px 0 12px;
  color: var(--text-secondary);
}

ul {
  margin: 0;
  padding-left: 18px;
  color: var(--text-secondary);
}
</style>
