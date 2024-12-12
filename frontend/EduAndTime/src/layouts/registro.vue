<script setup lang="ts">
// COMPLETED
import { ref, computed } from 'vue';
import BotonContrasteRosa from '../components/BotonContrasteRosa.vue';
import { useRouter } from 'vue-router';
import { authInstance } from '../utils/axios';
import { useToast } from 'vue-toastification';
import InputTexto from '../components/InputTexto.vue';
import TituloEAT from '../components/TituloEAT.vue';


const router = useRouter();
const toast = useToast();

// Declaración de propiedades reactivas
const nombre = ref('');
const correo = ref('');
const contrasena = ref('');
const reContrasena = ref('');
const contrasenasCoinciden = computed(() => contrasena.value === reContrasena.value);

// Método para registrar el usuario
async function registrarUsuario() {
  // Validaciones iniciales
  if (!nombre.value.trim() || !correo.value.trim() || !contrasena.value.trim() || !reContrasena.value.trim()) {
    toast.error('Por favor llena todos los campos.');
    return;
  }

  if (!contrasenasCoinciden.value) {
    toast.error('Las contraseñas deben ser iguales.');
    return;
  }

  try {
    const mutation = `
      mutation crearUsuario($input: UsuarioInput!) {
        crearUsuario(input: $input) {
          Exito
        }
      }
    `;

    // Construir el input dinámicamente con los datos disponibles
    const input = {
      Nombre: nombre.value.trim(),
      Correo: correo.value.trim(),
      Contrasena: contrasena.value.trim(),
      // Eliminar Imagen si no es necesario enviar el campo
    };

    // Preparar las variables para la mutación
    const variables = { input };

    const response = await authInstance.post('/query', { mutation, variables });

    // Validar respuesta
    const exito = response?.data?.data?.crearUsuario;

    if (exito) {
      toast.success('Registro completo. ¡Bienvenido!');
      router.push({ name: 'sesion' });
    } else {
      toast.error('Inténtalo nuevamente. Ocurrió un error.');
      console.log('Respuesta del servidor:', response);
    }
  } catch (error) {
    toast.error('Error al registrar el usuario. Por favor inténtalo más tarde.');
    console.error('Error al registrar el usuario:', error);
  }
}
</script>


<template>
  <div class="w-full h-full bg-azulOscuroEAT md:flex md:flex-row font-gabriela">

    <img src="../img/inicio/RectangleLibros.png" alt="" class="absolute z-0 w-full clip-diagonal md:clip-complete md:static md:w-1/2">

    <div class="flex relative z-10 flex-col gap-4 justify-end items-center w-full h-full md:gap-10">

      <TituloEAT Text="Registrate"/>

      <InputTexto v-model="nombre" type="text" placeholder="Tu nombre" required/>

      <InputTexto v-model="correo" type="email" placeholder="Tu correo" required/>

      <InputTexto v-model="contrasena" placeholder="Tu contraseña" type="password" required/>

      <InputTexto v-model="reContrasena" placeholder="Repite tu contraseña" type="password" required/>

      <BotonContrasteRosa label="Registrarse" @click="registrarUsuario"/>

      <p class="text-pink-700 font-abhayaRegular md:text-xl" v-if="!contrasenasCoinciden" >Las contraseñas no coinciden</p>

      <div class="flex flex-row gap-2 pl-20 mb-20 text-white md:mb-40 md:pl-60 md:text-xl">
        <p class="font-abhayaRegular">¿Ya tienes cuenta?</p>
        <router-link to="/" class="font-abhayaBold">Inicia sesión</router-link>
      </div>
    </div>
    
  </div>
</template>

