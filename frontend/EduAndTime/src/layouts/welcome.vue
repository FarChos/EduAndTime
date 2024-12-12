<script setup lang="ts">
//COMPLETED
  import { ref } from 'vue';
  const userData = JSON.parse(localStorage.getItem('userDataEAT') || '{}');
  const userName = ref(userData.Nombre || 'Usuario');
</script>
<template>
  <div class="flex flex-col gap-4 justify-center items-center w-full h-full md:flex-row md:gap-2 md:rounded-2xl bg-grisClaroEAT">
    <img class="h-60 md:h-40" src="../img/menu/logoAzulEAT.png" alt=""> 
    <div>
      <h1 class="text-4xl text-center md:text-7xl text-azulOscuroEAT md:text-start">Bienvenido {{ userName }}</h1>
      <h1 class="text-4xl text-center md:text-7xl text-azulOscuroEAT md:text-start">a Edu And Time</h1>
    </div>
  </div> 
</template> 