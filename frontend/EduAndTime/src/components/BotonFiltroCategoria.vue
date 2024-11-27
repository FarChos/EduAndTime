<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';

// Definir las propiedades esperadas
const props = defineProps({
  botonFocus: {
    type: String as () => string | null,
    required: false,
  },
  nombre: {
    type: String,
    required: true, // Asegúrate de que siempre tenga un valor válido
  },
  label: {
    type: String,
    required: true,
  },
});

// Emitir eventos al componente padre
const emit = defineEmits(['changeFocus']);

// Manejar el cambio de foco
const handleFocus = () => {
  emit('changeFocus', props.nombre); // Usar `props.nombre` explícitamente
};
</script>


<template>
  <button 
    :class="{
      'pr-1 pl-1 h-8 bg-azulOscuroEAT rounded-lg border-2 border-grisEAT text-white': botonFocus === nombre,
      'pr-1 pl-1 h-8 bg-white rounded-lg border-2 border-azulOscuroEAT': botonFocus !== nombre,
    }"
    @click="handleFocus"
  >
    {{ label }}
  </button>
</template>
