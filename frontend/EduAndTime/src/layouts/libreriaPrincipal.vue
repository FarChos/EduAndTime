<script setup lang="ts">
//COMPLETED
import { ref, onMounted } from 'vue';
import LibroMuestra from '../components/LibroMuestra.vue';
import BotonFiltroCategoria from '../components/BotonFiltroCategoria.vue';
import { useRouter } from 'vue-router';
import { useToast } from 'vue-toastification';
import { libreriaInstancia } from '../utils/axios';
import { formateadorEtiquetas } from '../utils/dataHelpers';
import { RecursoMuestra, GraphQLResponse } from '../utils/types';
import { useRutasGlobalRecursos } from '../stores/rutasGlobales';

const router = useRouter();
const toast = useToast();
const recomendadosMostrado = ref(false);

const parametroLiteral = ref('');
const Parametros = ref('');

const recursos = ref<RecursoMuestra[]>([]); // Inicializa el ref

const focusedButton = ref('general');
const rutasStore = useRutasGlobalRecursos()
const urlBase = rutasStore.rutaRecursos;

console.log(urlBase)

// Estado inicial del botón enfocado

const handleClick = (bookId: number) => {
  // Navegar a la ruta 'libroCompleto' pasando el id como parámetro
  router.push({ name: 'libroCompleto', query: { id: bookId } });
};
const irSubirRecurso= () => {
  router.push({ name: 'subirRecurso'});
}

const hacerPeticion = async (
  categoria: string, 
  cantidad: number, 
  titulo?: string, 
  autor?: string, 
  formato?: string, 
  etiquetas?: string[]
): Promise<GraphQLResponse<{ buscarRecursos: RecursoMuestra[] }>> => {
  try {
    
    const query = `
      query buscarRecursos($input: ParametrosBusqueda!) {
        buscarRecursos(input: $input) {
          Id
          Titulo
          Autor
          Categoria
          Formato
          Archivo
          Etiquetas
          Calificacion
        }
      }
    `;
    const input = {
      Titulo: titulo || null,
      Autor: autor || null,
      Categoria: categoria,
      Formato: formato || null,
      Etiquetas: etiquetas || null,
      Cantidad: cantidad
    };

    const variables = { input };
    console.log(input)
    // Realizar la consulta al backend
    const response = await libreriaInstancia.post('/query', { query, variables });
    console.log(response)
    // Verificar si hay errores en la respuesta
    if (response.data.errors) {
      const errorMessage = response.data.errors.map((err: any) => err.message).join(', ');
      throw new Error(errorMessage);
    }

    return response.data;
  } catch (error: any) {
    toast.error('Error al buscar recursos: ' + error.message);
    console.error('Error detallado:', error);
    return { errors: [{ message: 'Error al realizar la petición' }] };
  }
};
function asignarParametro(tipo: string, valor: string) {
  switch (tipo) {
    case 'titulo': return { titulo: valor };
    case 'autor': return { autor: valor };
    case 'etiquetas': return { etiquetas: formateadorEtiquetas(valor) };
    case 'formato': return { formato: valor };
    default: return {};
  }
}


const obtenerRecursos = async (cantidad: number): Promise<RecursoMuestra[] | undefined> => {
  recursos.value = []; // Reasigna un array vacío para "limpiar" los datos
  let parametros: { [key: string]: any } = { categoria: focusedButton.value, cantidad };
  if (parametroLiteral.value && Parametros.value) {
    parametros = { ...parametros, ...asignarParametro(Parametros.value, parametroLiteral.value) };
  }

  try {
    const respuesta = await hacerPeticion(
      parametros.categoria,
      parametros.cantidad,
      parametros.titulo,
      parametros.autor,
      parametros.formato,
      parametros.etiquetas
    );
    if (respuesta.errors) throw new Error(respuesta.errors.map(e => e.message).join(', '));
    console.log(respuesta.data?.buscarRecursos)
    return respuesta.data?.buscarRecursos || [];
  } catch (error) {
    console.error('Error al obtener recursos:', error);
    toast.error('No se pudo conectar con el servidor.');
    return undefined;
  }
};

const cargarInicialRecursos = async () : Promise<RecursoMuestra[] | undefined> => {
  
  try {
    const datos = await obtenerRecursos(10); // Obtener los recursos (9 en este caso)

    if (!datos) {
      toast.error('No se pudieron obtener los recursos.');
    } else {
      return datos; // Asignamos los recursos a la variable reactiva
    }
  } catch (error) {
    toast.error('Error al cargar los recursos.');
    console.error(error);
  }
};

// Manejar el cambio de foco en los botones
async function  buscar () {
  
  const resultadoFinal = await cargarInicialRecursos();
  
  // Asignar el resultado a recursos, si no es undefined
  if (resultadoFinal !== undefined) {
    recursos.value = resultadoFinal;  // Asigna el valor al ref
  }
}

// Manejar el cambio de foco en los botones
async function  cambiarFocus (button: string) {
  if (!button) {
    toast.error('Botón inválido.');
    return;
  }
  focusedButton.value = button; // Si lo haces reactivo
  const resultadoFinal = await cargarInicialRecursos();
  
  // Asignar el resultado a recursos, si no es undefined
  if (resultadoFinal !== undefined) {
    recursos.value = resultadoFinal;  // Asigna el valor al ref
  }
}

onMounted(async () => {
  
  // Llamar a cargarInicialRecursos y esperar el resultado
  const resultadoFinal = await cargarInicialRecursos();
  
  // Asignar el resultado a recursos, si no es undefined
  if (resultadoFinal !== undefined) {
    recursos.value = resultadoFinal;  // Asigna el valor al ref

 
  }
});

</script>

<template>
  <div class="flex flex-col w-full h-full bg-grisClaroEAT md:rounded-2xl">
    <header class="flex justify-center items-center w-full h-14 md:justify-start md:pl-2 md:h-16 md:pt-2 md:mb-1">
      <img
        class="absolute left-5 z-10 h-6 md:left-24"
        src="../img/libreria/iconoLupa.png"
        alt="Ícono de lupa"
      />

      <input
        @keydown.enter="buscar"
        v-model="parametroLiteral"
        :placeholder="`Buscar por ${Parametros || '...'}...`"
        class="pl-8 w-2/3 h-8 text-white border-2 rounded-s-xl focus:outline-none border-azulOscuroEAT bg-grisEAT placeholder:text-white md:w-3/6 md:h-12 md:text-lg"
        type="text"
      />

      <select 
      v-model="Parametros"
      class="pl-1 h-8 border-r-2 border-y-2 border-azulOscuroEAT rounded-e-xl md:h-12 md:text-lg text-azulOscuroEAT">
        <option selected disabled value=''>Parametros</option>
        <option value="titulo">Titulo</option>
        <option value="autor">Autor</option>
        <option value="etiquetas">Etiquetas</option>
        <option value="formato">Formato</option>
      </select>

    </header>
    <hr class="border-azulOscuroEAT">
    <nav class="flex overflow-x-auto gap-1 items-center pl-1 h-12 whitespace-nowrap scrollbar-hide md:justify-center">
      <BotonFiltroCategoria label="+" nombre="+" :botonFocus="focusedButton" @changeFocus="irSubirRecurso"/>
      <BotonFiltroCategoria label="General" nombre="general" :botonFocus="focusedButton" @changeFocus="cambiarFocus"/>
      <BotonFiltroCategoria label="Sociedad" nombre="sociedad" :botonFocus="focusedButton" @changeFocus="cambiarFocus"/>
      <BotonFiltroCategoria label="Geografia" nombre="geografia" :botonFocus="focusedButton" @changeFocus="cambiarFocus"/>
      <BotonFiltroCategoria label="Tecnologia" nombre="tecnologia" :botonFocus="focusedButton" @changeFocus="cambiarFocus"/>
      <BotonFiltroCategoria label="Ciencia" nombre="ciencia" :botonFocus="focusedButton" @changeFocus="cambiarFocus"/>
      <BotonFiltroCategoria label="Economia" nombre="economia" :botonFocus="focusedButton" @changeFocus="cambiarFocus"/>
      <BotonFiltroCategoria label="Bienestar" nombre="bienestar" :botonFocus="focusedButton" @changeFocus="cambiarFocus"/>
      <BotonFiltroCategoria label="Politica" nombre="politica" :botonFocus="focusedButton" @changeFocus="cambiarFocus"/>
      <BotonFiltroCategoria label="Arte" nombre="arte" :botonFocus="focusedButton" @changeFocus="cambiarFocus"/>
      <BotonFiltroCategoria label="Filosofia" nombre="filosofia" :botonFocus="focusedButton" @changeFocus="cambiarFocus"/>
      <BotonFiltroCategoria label="Exotica" nombre="exotica" :botonFocus="focusedButton" @changeFocus="cambiarFocus"/>
    </nav>
    
    <hr class="border-azulOscuroEAT">
    <main class="flex flex-col flex-grow items-center w-full h-full">
      <h1 class="md:text-xl text-azulOscuroEAT">Ultimas novedades</h1>
      <section class="w-full">
        <!-- Contenedor grid -->
        <div 
          class="grid overflow-y-auto lg:grid-cols-4  gap-4 max-h-[calc(3*theme(spacing.48)+2*theme(spacing.4))] grid-cols-1 ps-4 md:pe-4 "
        >
          <!-- Mostrar "Recomendados" solo una vez al inicio -->
          <div v-if="recursos.length >= 3 && !recomendadosMostrado" class="col-span-full">
            <h1 class="md:text-xl text-azulOscuroEAT">Recomendados</h1>
          </div>

          <div v-if="recursos.length === 0" class="text-center text-gray-500">
            No se encontraron recursos.
          </div>

          <!-- Mostrar cada recurso -->
          <LibroMuestra
            v-for="recurso in recursos"
            :key="recurso.Id"
            :epubUrl="urlBase + recurso.Archivo"
            :id="recurso.Id"
            :titulo="recurso.Titulo"
            :autor="recurso.Autor"
            :categoria="recurso.Categoria"
            :formato="recurso.Formato"
            :etiquetas="recurso.Etiquetas"
            :calificacion="recurso.Calificacion"
            @click="handleClick"
          />
        </div>
      </section>

    </main>
  </div>

</template>