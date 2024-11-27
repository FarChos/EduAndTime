<script setup lang="ts">
  import { ref } from 'vue';
  const userData = JSON.parse(localStorage.getItem('userDataEAT') || '{}');
  const userName = ref(userData.Nombre || 'Usuario');
</script>
<template>
  <div class="flex justify-center items-center w-full h-full md:rounded-2xl bg-grisClaroEAT">
    <img class="h-40" src="../img/menu/logoAzulEAT.png" alt=""> 
    <div>
      <h1 class="text-7xl text-azulOscuroEAT">Bienvenido {{ userName }}</h1>
      <h1 class="text-7xl text-azulOscuroEAT">a Edu And Time</h1>
    </div>
  </div> 
</template> 