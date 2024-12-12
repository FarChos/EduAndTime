<script setup lang="ts">
//COMPLETED
  import { useRouter } from 'vue-router'
  import { ref } from 'vue';
  import { useToast } from 'vue-toastification';

  const router = useRouter();
  const toast = useToast();
  const esVisinleMenu = ref(false)
  const focusedButton = ref(null);
  // Estado inicial del botón enfocado
  const rutaImagenCompleta = localStorage.getItem('userProfileImage');

  let imagenPerfil : string = '/src/img/menu/usuarioGenerico.jpg';
  
  if (rutaImagenCompleta != '/public/imagenesPerfil/' && rutaImagenCompleta != null){
    imagenPerfil = rutaImagenCompleta 
  }
  // Manejar el enfoque en los botones
  const mostrarMenu = () =>{
    esVisinleMenu.value = !esVisinleMenu.value
  }
  const irAMenu = () => {
    router.push({ name: 'menu' });
    focusedButton.value = null
  }
  const irACuenta = () => {
    router.push({ name: 'cuenta' });
    focusedButton.value = null
  }
  const cerrarSesion = () => {
  localStorage.removeItem('authTokenEAT'); // Elimina el token
  localStorage.removeItem('userDataEAT'); // Opcional: elimina datos de usuario
  router.push({ name: 'sesion' }); // Redirige al usuario a la página de inicio de sesión
  toast.success('Sesión cerrada exitosamente');
  };

  const handleFocus = (button: any) => {
    focusedButton.value = button;
    if(focusedButton.value === 'libreria'){
      router.push({ name: 'libreriaPrincipal' });
    }
    
  };
</script>

<template>
  <div class="flex flex-col w-full h-full font-gabriela">
    <!-- Header -->
    <header class="flex justify-center items-center p-1 w-full h-12 rounded-b-md bg-azulOscuroEAT md:rounded-none">
      <img @click="irAMenu"class="w-9 h-10 cursor-pointer md:ml-4" src="../img/menu/LogoEat.png" alt="Logo Eat">
      <h1 class="flex-grow text-lg text-center text-grisClaroEAT">Edu And Time</h1>
      <div class="relative group">
        <img
          @click="mostrarMenu"
          class="w-10 h-10 rounded-full border-2 border-transparent cursor-pointer active:border-grisEAT"
          :src='imagenPerfil'
          alt="Usuario Generico">
          <div name="menu" v-bind:class="{ hidden: !esVisinleMenu }" class="flex absolute right-10 flex-col gap-8 pt-6 w-60 h-40 rounded-2xl border-2 shadow pe-1 ps-1 bg-grisEAT border-azulOscuroEAT shadow-cyan-500/50">

            <button @click="irACuenta" class="h-10 text-xl rounded-xl border-2 shadow-sm border-azulOscuroEAT shadow-cyan-500/50 bg-grisClaroEAT text-azulOscuroEAT active:shadow-inner active:text-cyan-500/50 active:bg-gray-200 active:shadow-slate-500">Cuenta</button>
            <button @click="cerrarSesion" class="h-10 text-xl rounded-xl border-2 shadow-sm border-azulOscuroEAT shadow-cyan-500/50 bg-grisClaroEAT text-azulOscuroEAT active:shadow-inner active:text-cyan-500/50 active:bg-gray-200 active:shadow-slate-500">Cerrar sesión</button>

          </div>
      </div>
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
