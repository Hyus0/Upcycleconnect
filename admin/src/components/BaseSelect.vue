<template>
  <select :value="modelValue" @change="$emit('update:modelValue', $event.target.value)">
    <option v-for="option in options" :key="option.value" :value="option.value">
      {{ option.label }}
    </option>
  </select>
</template>

<script setup>
defineProps({
  modelValue: {
    type: [String, Number],
    default: ""
  },
  options: {
    type: Array,
    required: true
  }
});

defineEmits(["update:modelValue"]);
</script>
