<script setup lang="ts">
import { ref } from 'vue';
  import BarraBusqueda from '../components/BarraBusqueda.vue';
  import LibroMuestra from '../components/LibroMuestra.vue';
  import { useRouter } from 'vue-router'
  import { authInstance } from '../utils/axios';
  import { useToast } from 'vue-toastification';

const router = useRouter();
const toast = useToast();

// Declaración de propiedades reactivas
const correo = ref('');
const contrasena = ref('');
  async function iniciarSecion() {
  console.log('Iniciar sesión ejecutado');
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
    console.log(response.data);
    const autentificarUsuario = response.data.data.autentificarUsuario;
    const exito = autentificarUsuario.Exito;

    if (exito) {
      console.log('Usuario registrado:', response.data);

      // Guardar Token y datos del usuario en localStorage
      const { Token, Usuario } = autentificarUsuario;
      localStorage.setItem('authTokenEAT', Token); // Token
      localStorage.setItem('userDataEAT', JSON.stringify(Usuario)); // Otros datos del usuario

      // Redireccionar al home
      router.push({ name: 'menu' });
    } else {
      toast.info('Ups, revisa los campos, algo salió mal');
      console.log(response);
    }
  } catch (error) {
    toast.error('Error al iniciar sesión');
    console.error('Error al iniciar sesión:', error);
  }
}
</script>
<template>
  <div class="w-full h-full md:rounded-2xl bg-grisClaroEAT">
    <header class="flex justify-center items-center w-full h-14 md:justify-start md:pl-2 md:h-16">
      <BarraBusqueda nombreInput="" nombreSelec="" placeholder="Buscar. . ." :opciones="['Titulo','Autor','Etiquetas','Formato']" :valores="['titulo','autor','etiqueta','formato']"/>
    </header>
    <main class="scrollbar-hide ps-2 pe-2">
      <section class="w-full h-auto bg-grisEAT">
        <h1></h1>
      </section>
      <section>
        
      </section>
    </main>
  </div>
</template>