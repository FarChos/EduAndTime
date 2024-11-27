<script setup lang="ts">
import { ref } from 'vue';

// Estado inicial del botón enfocado
const focusedButton = ref(null);

// Manejar el enfoque en los botones
const handleFocus = (button: any) => {
  focusedButton.value = button;
};
</script>

<template>
  <div class="flex flex-col w-full h-full font-gabriela">
    <!-- Header -->
    <header class="flex justify-center items-center p-1 w-full h-12 rounded-b-md bg-azulOscuroEAT md:rounded-none">
      <img class="w-9 h-10 md:ml-4" src="../img/menu/LogoEat.png" alt="Logo Eat">
      <h1 class="flex-grow text-lg text-center text-grisClaroEAT">Edu And Time</h1>
      <img
        class="w-10 h-10 rounded-full border-2 border-transparent active:border-grisEAT"
        src="../img/menu/usuarioGenerico.jpg"
        alt="Usuario Generico"
      >
    </header>

    <!-- Contenedor principal -->
    <div class="flex flex-col w-full h-full md:flex-row-reverse md:w-full bg-azulOscuroEAT">
      <!-- Main -->
      <main class="flex-grow">
          <router-view></router-view> <!-- Aquí se cargará el contenido dinámico -->
      </main>

      <!-- Navegación -->
      <nav class="flex gap-4 justify-center items-center w-full h-8 bg-grisClaroEAT md:bg-azulOscuroEAT md:flex-col md:h-full md:w-20">
        <!-- Botón "Social" -->
        <button
        :class="{
            'w-2/5 rounded-md border-2 border-solid bg-grisClaroEAT text-azulOscuroEAT border-azulOscuroEAT md:border-none md:w-3/4 md:h-14': focusedButton === 'social',
            'w-2/5 rounded-md border-2 border-solid bg-azulOscuroEAT text-grisClaroEAT border-grisClaroEAT md:border-none md:w-3/4 md:h-14': focusedButton !== 'social',
          }"
          @click="handleFocus('social')"
        >
          Social
        </button>

        <!-- Botón "Libreria" -->
        <button
          :class="{
            'w-2/5 rounded-md border-2 border-solid bg-grisClaroEAT text-azulOscuroEAT border-azulOscuroEAT md:border-none md:w-3/4 md:h-14': focusedButton === 'libreria',
            'w-2/5 rounded-md border-2 border-solid bg-azulOscuroEAT text-grisClaroEAT border-grisClaroEAT md:border-none md:w-3/4 md:h-14': focusedButton !== 'libreria',
          }"
          @click="handleFocus('libreria')"
        >
          Libreria
        </button>
      </nav>
    </div>
  </div>
</template>
