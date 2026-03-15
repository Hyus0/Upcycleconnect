<template>
  <label class="field">
    <span class="label">{{ label }}</span>
    <slot />
    <small v-if="hint" class="hint">{{ hint }}</small>
    <small v-if="error" class="error">{{ error }}</small>
  </label>
</template>

<script setup>
defineProps({
  label: {
    type: String,
    required: true
  },
  hint: {
    type: String,
    default: ""
  },
  error: {
    type: String,
    default: ""
  }
});
</script>

<style scoped>
.field {
  display: grid;
  gap: 8px;
}

.label {
  font-weight: 700;
}

.field :deep(input),
.field :deep(select),
.field :deep(textarea) {
  width: 100%;
  min-height: 48px;
  padding: 12px 14px;
  border-radius: 14px;
  border: 1px solid var(--border);
  background: white;
  color: var(--text-primary);
}

.field :deep(textarea) {
  min-height: 120px;
  resize: vertical;
}

.hint {
  color: var(--text-muted);
}

.error {
  color: var(--accent-coral);
}
</style>
