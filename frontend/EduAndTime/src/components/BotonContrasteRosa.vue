<template>
  <button class="w-60 h-12 text-white bg-gradient-to-t rounded-xl active:shadow-2xl from-moradoClaro to-moradoOscuroEAT md:w-1/2 md:text-3xl md:h-20"
    @click="handleClick">
    {{ label }}
  </button>
</template>

<script setup lang="ts">
defineProps({
  label: {
    type: String,
    required: true,
  },
});

const handleClick = () => {
  // Aquí puedes llamar a la función directamente si quieres hacer algo más
  // o simplemente emular la lógica sin emitir el click.
};
</script>
