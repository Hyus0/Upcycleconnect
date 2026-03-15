<template>
  <input
    type="date"
    :value="modelValue"
    @input="$emit('update:modelValue', $event.target.value)"
  />
</template>

<script setup>
defineProps({
  modelValue: {
    type: String,
    default: ""
  }
});

defineEmits(["update:modelValue"]);
</script>
