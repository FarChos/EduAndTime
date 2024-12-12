<script setup lang="ts">
// COMPLETED
import { ref, computed } from 'vue';
import BotonContrasteRosa from '../components/BotonContrasteRosa.vue';
import { useRouter } from 'vue-router';
import { useToast } from 'vue-toastification';
import { authInstance } from '../utils/axios';

const router = useRouter();
const toast = useToast();

const reContrasena = ref('');
const contrasena = ref('');
const nombre = ref('');
const contrasenasCoinciden = computed(() => contrasena.value === reContrasena.value);
const containerRef = ref<HTMLDivElement | null>(null);
const imagen = ref<File | null>(null);
const fileName = ref<string | null>(null);

// Manejar el cambio del archivo
function handleFileChange(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (file) {
    fileName.value = file.name;
    imagen.value = file;

    // Renderizar la imagen seleccionada
    const imageUrl = URL.createObjectURL(file);
    const img = document.createElement('img');
    if (containerRef.value) {
      containerRef.value.innerHTML = '';
      img.src = imageUrl;
      img.alt = file.name;
      img.style.maxWidth = '50%';
      img.style.height = '100%';
      containerRef.value.appendChild(img);
    }

    // Liberar la URL después de usarla
    img.onload = () => URL.revokeObjectURL(imageUrl);
  }
}

async function actualizarPerfil() {
  // Validaciones de los campos de entrada
  if (!nombre.value.trim()) {
    toast.error('El nombre no puede estar vacío.');
    return;
  }
  if (!contrasena.value.trim()) {
    toast.error('La contraseña no puede estar vacía.');
    return;
  }
  if (!contrasenasCoinciden.value) {
    toast.error('Las contraseñas deben coincidir.');
    return;
  }
  if (!imagen.value) {
    toast.error('Por favor selecciona un archivo de imagen.');
    return;
  }

  // Obtener los datos del usuario y token de localStorage
  const userData = JSON.parse(localStorage.getItem('userDataEAT') || 'null');
  const userToken = JSON.parse(localStorage.getItem('authTokenEAT') || 'null');

  if (!userData || !userToken) {
    toast.error('Error de autenticación. Por favor inicia sesión.');
    return;
  }

  const correo = userData.Correo;
  const id = userData.Id;

  // Crear la mutación GraphQL
  const mutation = `
    mutation actualizarUsuario($id: Int!, $input: UsuarioInput!) {
      actualizarUsuario(id: $id, input: $input) {
        Exito
        Mensaje
      }
    }
  `;

  // Preparar el input
  const input: any = {
    Nombre: nombre.value.trim(),
    Correo: correo,
    Contrasena: contrasena.value.trim(),
  };

  // Si la imagen se ha seleccionado, incluirla en el input
  if (imagen.value) {
    input.Imagen = imagen.value;
  }

  // Preparar las variables para la mutación
  const variables = {
    id: id,
    input: input,
  };

  try {
    // Realizar la solicitud con GraphQL
    const response = await authInstance.post('/query', {
      query: mutation,
      variables,
    });

    const success = response.data.data.actualizarUsuario;

    if (success) {
      toast.success('Tus datos se actualizaron correctamente.');
      router.push({ name: 'cuenta' });
    } else {
      toast.error('Hubo un error al actualizar tus datos.');
      console.log('Error del servidor:', response.data.data);
    }
  } catch (error) {
    console.error('Error al actualizar:', error);
    toast.error('Algo salió mal al intentar actualizar tu perfil. Intenta más tarde.');
  }
}

</script>


<template>
  <main class="flex overflow-y-auto overflow-x-hidden flex-col gap-1 pt-2 bg-grisClaroEAT pe-2 ps-2 md:items-center md:justify-center md:gap-6">
    <article>
      <h1 class="text-xl font-bold text-center text-azulOscuroEAT">Actualiza tus datos</h1>
    </article>
    <article class="flex flex-col h-full rounded-xl bg-grisEAT pe-2 ps-2 md:ps-4 md:flex-row md:h-4/5 md:w-4/5 md:justify-between">
      <div class="flex flex-col md:w-1/2 md:justify-between">
        <div class="flex flex-col gap-3 md:gap-10 md:h-1/2">
          <h1 class="text-lg text-center text-azulProfundo md:text-2xl">Actualiza tu contraseña</h1>

          <input v-model="contrasena" placeholder="Nueva contraseña. . ." type="password" class="pl-2 w-full h-10 text-lg rounded-lg border-2 md:h-14 md:text-xl md:w-12/12 text-azulOscuroEAT border-azulOscuroEAT placeholder:text-azulGisaseoEAT bg-grisEAT focus:outline-none" required>

          <input v-model="reContrasena" placeholder="Repite la contraseña. . ." type="password" class="pl-2 w-full h-10 text-lg rounded-lg border-2 md:h-14 md:text-xl md:w-12/12 text-azulOscuroEAT border-azulOscuroEAT placeholder:text-azulGisaseoEAT bg-grisEAT focus:outline-none" required>

        </div>
        <div class="flex flex-col gap-3 mt-6 md:h-1/2">
          <h1 class="text-lg text-center text-azulProfundo md:text-2xl">Actualiza tu nombre</h1>

          <input v-model="nombre" placeholder="Escribe tu nombre. . ." type="text" class="pl-2 w-full h-10 text-lg rounded-lg border-2 md:h-14 md:text-xl md:w-12/12 text-azulOscuroEAT border-azulOscuroEAT placeholder:text-azulGisaseoEAT bg-grisEAT focus:outline-none" required>

        </div>
      </div>
      <div class="flex flex-col gap-3 items-center mt-6 md:w-1/2 md:mt-0">
        <h1 class="text-lg text-center text-azulProfundo md:text-2xl">Actualiza tu imagen de perfil</h1>
        <div class="flex flex-col gap-2 items-end w-full">
            <span v-if="fileName" class="ml-2 text-black">{{ fileName }}</span>
            <input
              type="file"
              id="file-input"
              class="hidden"
              accept=".png, .jpg"
              @change="handleFileChange"
              required
            />
            <!-- Label personalizado que actúa como botón -->
            <label for="file-input" class="p-2 w-1/2 text-center text-white rounded-lg shadow-xl cursor-pointer active:shadow-none bg-azulOscuroEAT">
              Subir imagen
            </label>
          </div>
          <BotonContrasteRosa class="text-xl md:text-xl md:h-9" label="Actualizar perfil" @click="actualizarPerfil"/>
          <div ref="containerRef" style="width: 75%; height: 150px; border: 1px solid #ddd;"></div>
      </div>
    </article>
  </main>
</template>