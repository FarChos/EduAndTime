<script setup lang="ts">
//COMPLETED
import { ref, nextTick } from 'vue';
import ePub from 'epubjs';
import { libreriaInstancia } from '../utils/axios';
import { useToast } from 'vue-toastification';
import { useRouter } from 'vue-router';
import { formateadorEtiquetas, formatoLibro } from '../utils/dataHelpers';

const router = useRouter();
const toast = useToast();

// Estado reactivo
const titulo = ref('');
const autor = ref('');
const descripcion = ref('');
const etiqueta = ref('');
const categoria = ref('');
const recurso = ref<File | null>(null);

const containerRef = ref<HTMLDivElement | null>(null);
const epubBlob = ref<ArrayBuffer | null>(null);
const fileName = ref('');

// Manejar la selección del archivo
const handleFileChange = (event: Event) => {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];

  if (!file) return;

  recurso.value = file; // Asigna el archivo seleccionado

  if (file.type === 'application/epub+zip') {
    fileName.value = file.name;

    const reader = new FileReader();
    reader.onload = async (e) => {
      if (e.target?.result && containerRef.value) {
        containerRef.value.innerHTML = ''; // Limpiar contenedor
        epubBlob.value = e.target.result as ArrayBuffer;
        await renderEpub(epubBlob.value); // Renderizar el archivo EPUB
      }
    };
    reader.onerror = () => {
      toast.error('Error al leer el archivo. Asegúrate de que esté en formato EPUB.');
    };
    reader.readAsArrayBuffer(file);
  } else {
    toast.error('Por favor selecciona un archivo EPUB válido.');
  }
};

// Renderizar el EPUB
const renderEpub = async (epubData: ArrayBuffer | null) => {
  try {
    if (containerRef.value && epubData) {
      await nextTick();
      const book = ePub(epubData);
      const rendition = book.renderTo(containerRef.value, {
        width: '100%',
        height: '130%',
      });

      await book.ready;
      rendition.themes.default({
        '*': {
          margin: '0 !important',
          padding: '0 !important',
        },
      });

      rendition.display();
    }
  } catch (error) {
    console.error('Error al renderizar el libro:', error);
  }
};

// Enviar el formulario
const submitForm = async () => {
  if (!titulo.value.trim() || !autor.value.trim() || !categoria.value.trim()) {
    toast.error('Por favor completa todos los campos obligatorios.');
    return;
  }

  if (!recurso.value) {
    toast.error('Por favor selecciona un archivo.');
    return;
  }

  
  const userData = JSON.parse(localStorage.getItem('userDataEAT') || 'null');

  if (!userData) {
    toast.error('Error de autenticación. Por favor inicia sesión.');
    return;
  }

  const idUsuario = userData.Id;
  const etiquetas = formateadorEtiquetas(etiqueta.value);
  const formato = formatoLibro(fileName.value);
  const formData = new FormData();
  
  // Crear la mutación y las variables
  //const archivo = document.querySelector('file-input').files[0]; 
  const archivo = recurso.value;
  formData.append('operations', JSON.stringify({
  query: `
    mutation subirRecurso($input: RecursoInput!) {
      subirRecurso(input: $input) {
        Exito
      }
    }
  `,
  variables: {
    input: {
      Titulo: titulo.value.trim(),
      Autor: autor.value.trim(),
      Categoria: categoria.value.trim(),
      IdUsuario: idUsuario,
      Formato: formato,
      Descripcion: descripcion.value.trim(),
      Recurso: null, // El archivo será asignado en el map
      Etiquetas: etiquetas,
    },
  },
}));

  formData.append('map', JSON.stringify({ '1': ['variables.input.Recurso'] }));

  if (archivo) {
    formData.append('1', archivo);
  } else {
    toast.error('Error: No se seleccionó un archivo válido.');
    return;
  }

  // Preparar las variables para la mutación
  try {
    const response = await libreriaInstancia.post('/query', formData, {
      headers: {
        'Content-Type': 'multipart/form-data' // Axios manejará el boundary automáticamente
      }
    });
    console.log(response.data);
    if (response.data.data.subirRecurso.Exito) {
      toast.success('Recurso subido correctamente.');
      router.push({ name: 'libreriaPrincipal' });
    } else {
      toast.error('Hubo un error al subir el recurso.');
    }

  } catch (error) {
    console.error('Error al subir el recurso:', error);
    toast.error('Algo salió mal al subir el recurso.');
  }
};
</script>




<template>
  <div class="flex justify-center items-center w-full h-full md:rounded-2xl bg-grisClaroEAT">
    <!-- Contenedor con scroll vertical y sin scroll horizontal -->
    <main class="overflow-y-auto overflow-x-hidden w-11/12 h-auto rounded-3xl shadow-xl md:h-5/6 md:w-5/6 bg-grisEAT ps-2 pe-2">
      <article class="flex justify-center items-center h-6 md:justify-start md:ps-6 md:h-14">
        <h1 class="text-lg font-bold text-azulOscuroEAT">{{ titulo }}</h1>
      </article>
      <hr class="border-azulAunMasOscuroEAT">
      <article class="text-lg md:ps-6">
        <div class="flex flex-col w-full md:flex-row md:mt-10">
          <div class="flex flex-col gap-4 md:w-3/5 md:gap-6">
            <div class="flex flex-col gap-2">
              <h1 class="text-azulOscuroEAT">Detalles</h1>
              <input type="text" class="pl-2 w-full rounded-lg border-2 md:h-11 md:w-11/12 text-azulOscuroEAT border-azulOscuroEAT placeholder:text-azulGisaseoEAT bg-grisEAT focus:outline-none" placeholder="Título. . ." v-model="titulo" required />
            </div>  
            <input type="text" class="pl-2 w-full rounded-lg border-2 md:h-11 md:w-11/12 text-azulOscuroEAT border-azulOscuroEAT placeholder:text-azulGisaseoEAT bg-grisEAT focus:outline-none" placeholder="Autor. . ." required v-model="autor" />

            <textarea class="pl-2 rounded-lg border-2 text-azulOscuroEAT md:w-11/12 md:h-60 border-azulOscuroEAT placeholder:text-azulGisaseoEAT bg-grisEAT focus:outline-none" name="" id="" placeholder="Descripción. . . " required v-model="descripcion"></textarea>
            <div class="flex flex-col">
              <label for="etiquetas">Las etiquetas inician con un '#' y terminan con espacio ' '</label>
              <input class="pl-2 mb-4 rounded-lg border-2 md:w-11/12 md:h-8 text-azulOscuroEAT border-azulOscuroEAT placeholder:text-azulGisaseoEAT bg-grisEAT focus:outline-none" type="text" id="etiquetas" placeholder="Etiquetas. . ." v-model="etiqueta" required />

              <select class="pl-2 rounded-lg border-2 text-azulOscuroEAT md:h-8 md:w-11/12 border-azulOscuroEAT placeholder:text-azulGisaseoEAT bg-grisEAT focus:outline-none" name="categoria" id="categoria" required v-model="categoria">
                
                <option value="" disabled selected>Categoría</option>
                <option value="sociedad">Sociedad</option>
                <option value="geografia">Geografía</option>
                <option value="tecnologia">Tecnología</option>
                <option value="ciencia">Ciencia</option>
                <option value="economia">Economía</option>
                <option value="bienestar">Bienestar</option>
                <option value="politica">Política</option>
                <option value="arte">Arte</option>
                <option value="filosofia">Filosofía</option>
                <option value="exotica">Exótica</option>
              </select>
          </div>
          </div>
        
        <div class="flex flex-col flex-grow gap-2 justify-end pb-4 mt-2 md:flex-col-reverse pe-6">
          <div class="flex flex-col gap-2 items-end">
            <span v-if="fileName" class="ml-2 text-black">{{ fileName }}</span>
            <input
              type="file"
              id="file-input"
              class="hidden"
              accept=".pdf, .epub"
              @change="handleFileChange"
            />
            <!-- Label personalizado que actúa como botón -->
            <label for="file-input" class="p-2 w-1/2 text-center text-white rounded-lg shadow-xl cursor-pointer active:shadow-none bg-azulOscuroEAT">
              Subir archivo
            </label>
            <button class="w-60 h-12 text-white rounded-lg shadow-xl active:shadow-none active:bg-moradoClaro bg-moradoOscuroEAT md:w-full md:text-3xl md:h-14"
              @click="submitForm" type="submit">
              Subir recurso
            </button>
          </div>
          <div ref="containerRef" style="width: 75%; height: 200px; border: 1px solid #ddd;"></div>
        </div>
      </div>
      </article>
    </main>
  </div>
</template>



