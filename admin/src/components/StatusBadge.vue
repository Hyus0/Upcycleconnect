<template>
  <span class="badge" :class="`badge-${tone}`">{{ label }}</span>
</template>

<script setup>
defineProps({
  label: {
    type: String,
    required: true
  },
  tone: {
    type: String,
    default: "green"
  }
});
</script>

<style scoped>
.badge {
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  padding: 7px 10px;
  font-size: 0.8rem;
  font-weight: 700;
}

.badge-green {
  background: var(--brand-green-pale);
  color: var(--brand-green);
}

.badge-amber {
  background: rgba(232, 160, 48, 0.16);
  color: #9a6316;
}

.badge-coral {
  background: rgba(217, 95, 75, 0.16);
  color: var(--accent-coral);
}

.badge-teal {
  background: rgba(27, 138, 138, 0.14);
  color: var(--accent-teal);
}

.badge-sand {
  background: rgba(242, 232, 217, 0.9);
  color: var(--brand-earth);
}
</style>
