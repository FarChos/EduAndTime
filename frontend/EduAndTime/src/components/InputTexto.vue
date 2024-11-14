<template>
  <input
    :placeholder="placeholder"
    :type="type" 
    :required="required"
    :value="modelValue" 
    @input="onInput"
    class="pl-4 w-80 h-12 text-white bg-transparent border-b-4 border-azulClaroEAT focus:outline-none focus:border-grisEAT md:w-1/2 md:text-2xl"
  >
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';

defineProps({
  type: {
    type: String,
    default: 'text',
  },
  placeholder: {
    type: String,
    required: true,
  },
  modelValue: String,
  required: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['update:modelValue']);

function onInput(event: Event) {
  const target = event.target as HTMLInputElement;
  emit('update:modelValue', target.value);
}
</script>
