<script setup lang="ts">
import { ref } from 'vue';
import BotonContrasteRosa from '../components/BotonContrasteRosa.vue';
import InputTexto from '../components/InputTexto.vue';
import TituloEAT from '../components/TituloEAT.vue';
import { useRouter } from 'vue-router'
import { authInstance } from '../utils/axios';
import { useToast } from 'vue-toastification';


const router = useRouter();
const toast = useToast();

// Declaración de propiedades reactivas
const correo = ref('');
const contrasena = ref('');

async function iniciarSecion() {
  try {
    const response = await authInstance.get('/query', {
      params: {
        query: `
        query {
          autentificarUsuario(correo: "${correo.value}", contrasena: "${contrasena.value}") {
            Token
            Usuario {
              Nombre
              Correo
              ImgPerf
            }
            Exito
            Expira
          }
        }
      `
      },
    });
    console.log(response.data)
    const exito = response.data.data.autentificarUsuario.Exito;

    if (exito){
      console.log('Usuario registrado:', response.data);
      const token = response.data.data.autentificarUsuario.Token;
      localStorage.setItem('authTokenEAT', token);
      router.push({ name: 'home' });
    }else {
      toast.info('Ups, revisa los campos, algo salio mal')
      console.log(response)
    }

  } catch (error) {
    toast.error('Error al iniciar sesion')
    console.error('Error al iniciar sesion:', error);
    console.log()
  }
}
</script>

<template>
  <div class="w-full h-full bg-azulOscuroEAT md:flex md:flex-row font-gabriela">

    <img src="../img/inicio/RectangleLibros.png" alt="" class="absolute z-0 w-full clip-diagonal md:clip-complete md:static md:w-1/2">

    <div class="flex relative z-10 flex-col gap-6 justify-end items-center w-full h-full md:gap-20">

      <TituloEAT Text="Inicio de sesión"/>
      <InputTexto v-model="correo" type="email" placeholder="Tu correo"/>

      <InputTexto v-model="contrasena" placeholder="Tu contraseña" type="password"/>

      <BotonContrasteRosa label="Iniciar sesión" @click="iniciarSecion"/>
      
      <div class="flex flex-row gap-2 pl-20 mb-20 text-white md:mb-40 md:pl-60 md:text-xl">
        <p class="font-abhayaRegular">¿No tienes cuenta?</p>
        <router-link to="/registro" class="font-abhayaBold">Registrate</router-link>
      </div>
    </div>
    
  </div>
</template>

