<template>
  <button class="w-60 h-12 text-white bg-gradient-to-t rounded-xl active:shadow-2xl from-moradoClaro to-moradoOscuroEAT md:w-1/2 md:text-3xl md:h-20"
    @click="$emit('click')">
    {{ label }}
  </button>
</template>

<script setup lang="ts">
defineProps({
  label: {
    type: String,
    required: true,
  },
  handleClick: {
    type: Function,
    default: () => {},
  }
});
</script>