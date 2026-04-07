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
  color: #f1f6f3;
}

.field :deep(input),
.field :deep(select),
.field :deep(textarea) {
  width: 100%;
  min-height: 48px;
  padding: 12px 14px;
  border-radius: 14px;
  border: 1px solid rgba(140, 177, 156, 0.18);
  background:
    linear-gradient(180deg, rgba(20, 31, 27, 0.96), rgba(13, 22, 18, 0.98));
  color: var(--text-primary);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.04),
    0 0 0 1px rgba(0, 0, 0, 0.12);
}

.field :deep(textarea) {
  min-height: 120px;
  resize: vertical;
}

.field :deep(input::placeholder),
.field :deep(textarea::placeholder) {
  color: rgba(167, 187, 176, 0.78);
}

.field :deep(select) {
  color-scheme: dark;
}

.hint {
  color: var(--text-muted);
}

.error {
  color: var(--accent-coral);
}
</style>
