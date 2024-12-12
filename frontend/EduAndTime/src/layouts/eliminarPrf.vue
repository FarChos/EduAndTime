<script setup lang="ts">
// COMPLETED
import BotonContrasteRosa from '../components/BotonContrasteRosa.vue';
import InputTexto from '../components/InputTexto.vue';
import { authInstance } from '../utils/axios';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from 'vue-toastification';

const toast = useToast();
const router = useRouter();

const contrasena = ref('');
const userData = JSON.parse(localStorage.getItem('userDataEAT') || '{}');
const correo = ref(userData.Correo);

async function EliminarUsuario() {
  if (!correo.value || !contrasena.value) {
    toast.error('Por favor, completa todos los campos');
    return;
  }

  const query = `
  mutation EliminarUsuario($input: UsuarioInput!) {
    eliminarUsuario(input: $input) {
      Exito
      Mensaje
    }
  }
`;

// Construir el input dinámicamente con los datos necesarios
const input = {
  Correo: correo.value.trim(),
  Contrasena: contrasena.value.trim(),
};

const variables = { input };

try {
  const response = await authInstance.post('/query', { query, variables });

  // Validar respuesta del servidor
  const eliminarUsuario = response?.data?.data?.eliminarUsuario;

  if (!eliminarUsuario) {
    toast.error('Error inesperado al intentar eliminar el usuario.');
    console.error('Respuesta inesperada:', response);
    return;
  }

  const { Exito, Mensaje } = eliminarUsuario;

  if (Exito) {
    toast.success(Mensaje || 'Usuario eliminado correctamente.');
    router.push({ name: 'inicio' }); // Ajusta la ruta según tu app
  } else {
    toast.error(Mensaje || 'No se pudo eliminar el usuario. Verifica los datos e inténtalo de nuevo.');
  }
} catch (error) {
  console.error('Error al eliminar el usuario:', error);
  toast.error('Ocurrió un error inesperado. Por favor, inténtalo más tarde.');
}

}
</script>

<template>
  <main class="flex flex-col w-full bg-azulOscuroEAT">
    <article class="flex gap-2 justify-center items-center">
      <img class="h-10" src="../img/inicio/LogoEat.png" alt="">
      <h1 class="text-xl text-white md:text-4xl">Eliminar Cuenta</h1>
    </article>
    <article class="flex flex-col flex-grow gap-6 justify-center items-center md:gap-40">
      <h1 class="text-lg text-center text-white md:text-3xl">Introduce tu contraseña para eliminar tu cuenta</h1>
      <InputTexto class="md:w-1/3" v-model="contrasena" type="password" placeholder="Tu contraseña" required/>
      <BotonContrasteRosa class="md:w-1/3 md:h-10" label="Eliminar cuenta" @click="EliminarUsuario()"/>
    </article>
  </main>
</template>