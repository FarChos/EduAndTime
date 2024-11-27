<script setup lang="ts">
const { placeholder, opciones, valores, nombreInput, nombreSelec } = defineProps<{
  nombreInput: string;
  nombreSelec: string;
  placeholder: string;
  opciones: string[];
  valores: string[];
}>();
</script>

<template>
  <img
    class="absolute left-5 z-10 h-6 md:left-24"
    src="../img/libreria/iconoLupa.png"
    alt="Ícono de lupa"
  />

  <input
    class="pl-8 w-2/3 h-8 text-white border-2 rounded-s-xl focus:outline-none border-azulOscuroEAT bg-grisEAT placeholder:text-white md:w-3/6 md:h-12 md:text-lg"
    :name="nombreInput"
    :placeholder="placeholder"
    type="text"
  />

  <select
    class="pl-1 h-8 border-r-2 border-y-2 border-azulOscuroEAT rounded-e-xl md:h-12 md:text-lg"
    :name="nombreSelec"
  >
    <option
      v-for="(opcion, index) in opciones"
      :key="index"
      :value="valores[index]"
    >
      {{ opcion }}
    </option>
  </select>
</template>
