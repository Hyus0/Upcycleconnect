<template>
  <div v-if="open" class="backdrop" @click.self="$emit('cancel')">
    <div class="modal surface-card">
      <h3>{{ title }}</h3>
      <p>{{ message }}</p>
      <div class="toolbar">
        <button class="button button-secondary" @click="$emit('cancel')">Annuler</button>
        <button class="button button-danger" @click="$emit('confirm')">{{ confirmLabel }}</button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  open: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: "Confirmer l'action"
  },
  message: {
    type: String,
    default: "Cette action est irreversible."
  },
  confirmLabel: {
    type: String,
    default: "Confirmer"
  }
});

defineEmits(["cancel", "confirm"]);
</script>

<style scoped>
.backdrop {
  position: fixed;
  inset: 0;
  background: rgba(26, 31, 28, 0.45);
  display: grid;
  place-items: center;
  padding: 18px;
  z-index: 900;
}

.modal {
  width: min(480px, 100%);
  padding: 24px;
}

.modal p {
  margin: 8px 0 18px;
  color: var(--text-secondary);
}
</style>
