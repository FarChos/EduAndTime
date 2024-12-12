<script setup lang="ts">
//COMPLETED
import { ref } from 'vue';
import BotonContrasteRosa from '../components/BotonContrasteRosa.vue';
import InputTexto from '../components/InputTexto.vue';
import TituloEAT from '../components/TituloEAT.vue';
import { useRouter } from 'vue-router';
import { authInstance } from '../utils/axios';
import { useToast } from 'vue-toastification';
import { useRutasGlobalImg } from '../stores/rutasGlobales';

const rutasStore = useRutasGlobalImg()
const router = useRouter();
const toast = useToast();

// Declaración de propiedades reactivas
const correo = ref('');
const contrasena = ref('');

async function iniciarSesion() {
  console.log('Iniciar sesión ejecutado');
  
  // Validar campos
  if (!correo.value.trim() || !contrasena.value.trim()) {
    toast.error('Por favor llena todos los campos');
    return;
  }

  try {
    const query = `
      query autentificarUsuario($input: UsuarioInput!) {
        autentificarUsuario(input: $input) {
          Token
          Usuario {
            Id
            Nombre
            Correo
            nombreImagen
          }
          Exito
        }
      }
    `;

    const input = {
      Correo: correo.value.trim(),
      Contrasena: contrasena.value.trim(),
    };

    const variables = { input };

    const response = await authInstance.post('/query', { query, variables });

    // Validar respuesta del servidor
    const autentificarUsuario = response?.data?.data?.autentificarUsuario;
    
    if (!autentificarUsuario) {
      toast.error('Error inesperado, por favor inténtalo más tarde.');
      console.error('Respuesta inesperada:', response);
      return;
    }

    const exito = autentificarUsuario.Exito;
    if (exito) {
      console.log('Usuario registrado:', response.data);

      // Guardar Token y datos del usuario en localStorage
      const rutaImagen = rutasStore.rutaImagenesPerfil;
      const Token = autentificarUsuario.Token;
      const Usuario = autentificarUsuario.Usuario;
      
      let rutaImagenCompleta = ''
      // Concatenar la ruta global con el nombre de la imagen del usuario
      if(Usuario.nombreImagen != null){
        rutaImagenCompleta = `${rutaImagen}${Usuario.nombreImagen}`;
      }

      // Guardar la ruta de la imagen y otros datos en localStorage
      localStorage.setItem('authTokenEAT', `Bearer ${Token}`); // Guardar token
      localStorage.setItem('userDataEAT', JSON.stringify(Usuario)); // Otros datos del usuario
      localStorage.setItem('userProfileImage', rutaImagenCompleta); // Guardar la ruta completa de la imagen


      toast.success('¡Bienvenido!');
      router.push({ name: 'menu' }); // Redirigir al menú principal
    } else {
      toast.info('Revisa los campos, algo salió mal');
      console.log(response);
    }
  } catch (error) {
    // Manejo de errores de red u otros
    toast.error('Error al iniciar sesión');
    console.error('Error al iniciar sesión:', error);
  }
}
</script>


<template>
  <div class="w-full h-full bg-azulOscuroEAT md:flex md:flex-row font-gabriela">

    <img src="../img/inicio/RectangleLibros.png" alt="" class="absolute z-0 w-full clip-diagonal md:clip-complete md:static md:w-1/2">
    <img class="" alt="">
    <div class="flex relative z-10 flex-col gap-6 justify-end items-center w-full h-full md:gap-20">

      <TituloEAT Text="Inicio de sesión"/>
      <InputTexto v-model="correo" type="email" placeholder="Tu correo"/>

      <InputTexto v-model="contrasena" placeholder="Tu contraseña" type="password"/>

      <BotonContrasteRosa label="Iniciar sesión" @click="iniciarSesion"/>
      <div class="flex flex-row gap-2 pl-20 mb-20 text-white md:mb-40 md:pl-60 md:text-xl">
        <p class="font-abhayaRegular">¿No tienes cuenta?</p>
        <router-link to="/registro" class="font-abhayaBold">Registrate</router-link>
      </div>
    </div>
    
  </div>
</template>

