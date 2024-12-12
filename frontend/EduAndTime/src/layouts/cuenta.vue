<script setup lang="ts">
//COMPLETED
  import { useRouter } from 'vue-router'
  import { ref } from 'vue';
  
  let esVisinleOpciones = ref(true);
  let esVisibleElContenido = ref(false);
  let volver = ref(true)
  const isSmallScreen = window.matchMedia('(max-width: 768px)').matches; // Define el tamaño según tus breakpoints
  if(!isSmallScreen){
    esVisibleElContenido = ref(true)
    volver = ref(false)
  }
  console.log (isSmallScreen)
  const router = useRouter();
  const userData = JSON.parse(localStorage.getItem('userDataEAT') || '{}');
  const userName = userData.Nombre || 'Usuario Generico';
  const rutaImagenCompleta = localStorage.getItem('userProfileImage');

  let imagenPerfil : string = '/src/img/menu/usuarioGenerico.jpg';
  if (rutaImagenCompleta != '' && rutaImagenCompleta != null){
    imagenPerfil = rutaImagenCompleta 
  }
  console.log(imagenPerfil)
  console.log(userName)

  const activarContenido = () =>{
    esVisinleOpciones = ref(false)
    esVisibleElContenido = ref(true)
  }
  const irAFavoritos = () => {
    if (isSmallScreen) {
      activarContenido();
    }
    router.push({ name: 'favoritos' });
  };
  const irAMisRecursos = () =>{
    if (isSmallScreen) {
      activarContenido();
    }
    router.push({ name: 'misRecursos' });
  }
  const irAEdicionPerfil = () =>{
    if (isSmallScreen) {
      activarContenido();
    }
    router.push({ name: 'editarPerfil' });
  }
  const irAEliminarPerfil = () =>{
    activarContenido()
    volver = ref(true)
    router.push({ name: 'eliminarPerfil' });
  }
  const volverAOpciones = () =>{
    esVisinleOpciones = ref(true)
    if(isSmallScreen){
      esVisibleElContenido = ref(false)
    }else{
      volver = ref(false)
    }
    
    router.push({ name: 'cuenta' });
  }
  const irAMenu = () => {
    router.push({ name: 'menu' });
  }
  
</script>
<template>
  <main class="flex w-full h-full flex-c font-gabriela bg-grisClaroEAT md:flex-row">
    <section v-bind:class="{ hidden: !esVisinleOpciones }"  class="flex flex-col w-full h-full md:w-1/5">
      <article class="flex flex-col gap-4 items-center pt-6 w-full shadow bg-azulOscuroEAT shadow-azulOscuroEAT ps-2 pe-2 md:rounded-e-xl">
        <div class="inline-block relative group">
            <img @click="irAMenu" class="w-12 cursor-pointer" src="../img/inicio/LogoEat.png" alt="Logo de edu and time">
            <div
            class="absolute left-1/2 px-4 py-2 mb-2 text-sm text-white rounded-md opacity-0 transition-opacity duration-300 transform -translate-x-1/2 bg-azulGisaseoEAT bottom-30 group-hover:opacity-100">
              <p>home</p>
            </div>
        </div>
        <div class="flex gap-2 justify-center items-center p-1 rounded-xl shadow">
          <img :src='imagenPerfil' alt="" class="w-14 h-14 rounded-full md:h-10 md:w-10">
          <h1 class="text-3xl text-center text-grisClaroEAT md:text-xl">{{ userName }}</h1>
        </div>
      </article>
      <article class="flex flex-col flex-grow gap-8 pt-8 pe-3 ps-3">
        <hr class="border-azulOscuroEAT">
        <button @click="irAFavoritos" class="justify-center items-center h-16 text-xl rounded-lg shadow cursor-pointer text-azulOscuroEAT shadow-azulOscuroEAT bg-grisEAT active:shadow-inner active:shadow-azulOscuroEAT">Favoritos</button>
        <hr class="border-azulOscuroEAT">

        <button @click="irAMisRecursos" class="justify-center items-center h-16 text-xl rounded-lg shadow cursor-pointer text-azulOscuroEAT shadow-azulOscuroEAT bg-grisEAT active:shadow-inner active:shadow-azulOscuroEAT">Tus recursos</button>
        <hr class="border-azulOscuroEAT">

        <button @click="irAEdicionPerfil" class="justify-center items-center h-16 text-xl rounded-lg shadow cursor-pointer text-azulOscuroEAT shadow-azulOscuroEAT bg-grisEAT active:shadow-inner active:shadow-azulOscuroEAT">Editar perfil</button>
        <hr class="border-azulOscuroEAT">

        <button @click="irAEliminarPerfil" class="justify-center items-center h-16 text-xl rounded-lg shadow cursor-pointer active:shadow-rojoEAT text-rojoEAT shadow-zinc-600 bg-grisEAT active:shadow-inner">Eliminar perfil</button>
        <hr class="border-azulOscuroEAT">
      </article>
    </section>
    <section class="flex flex-col w-full h-full md:3/5" v-bind:class="{ hidden: !esVisibleElContenido }">
      <div class="bg-azulOscuroEAT pe-2"
      v-bind:class="{ hidden: !volver }">
        <div class="p-1 w-12 rounded-lg active:bg-white">
          <img @click="volverAOpciones" class="w-10 h-10 cursor-pointer" src="../img/iconosLibro/iconoVolverBlanco.png" alt="flecha volver">
        </div>
      </div>
      <router-view class="flex-grow"></router-view>
    </section>
  </main>
</template>