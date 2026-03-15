<template>
  <article class="surface-card stat-card">
    <div class="stat-head">
      <span>{{ label }}</span>
      <StatusBadge :label="badgeLabel" :tone="tone" />
    </div>
    <strong class="stat-value">{{ valueLabel }}</strong>
    <p>{{ caption }}</p>
  </article>
</template>

<script setup>
import { computed } from "vue";
import { formatNumber } from "../utils/format";
import StatusBadge from "./StatusBadge.vue";

const props = defineProps({
  label: {
    type: String,
    required: true
  },
  value: {
    type: [Number, String, null],
    default: null
  },
  tone: {
    type: String,
    default: "green"
  },
  missing: {
    type: Boolean,
    default: false
  }
});

const valueLabel = computed(() => (props.missing ? "API manquante" : formatNumber(props.value)));
const badgeLabel = computed(() => (props.missing ? "a completer" : "connecte"));
const caption = computed(() =>
  props.missing
    ? "Le backend doit encore exposer cette ressource pour la supervision admin."
    : "Valeur consolidee depuis les endpoints disponibles."
);
</script>

<style scoped>
.stat-card {
  padding: 18px;
}

.stat-head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  color: var(--text-secondary);
}

.stat-value {
  display: block;
  margin: 16px 0 10px;
  font-family: "Syne", sans-serif;
  font-size: 2.2rem;
  line-height: 0.95;
}

p {
  margin: 0;
  color: var(--text-secondary);
}
</style>
