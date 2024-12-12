<script setup lang="ts">
  // Importaciones de Vue y librerías
  import { ref, onMounted, nextTick } from 'vue';
  import { useToast } from 'vue-toastification';
  import ePub from 'epubjs';
  
  // Importaciones de tipos y componentes
  import { Recurso, GraphQLResponse } from '../utils/types';
  import BarraBusqueda from '../components/BarraBusqueda.vue';
  

  // Importación de imágenes y utilidades
  import GustoAltoIcono from '../img/iconosLibro/GustoAlto.png';
  import GustoAltoMedio from '../img/iconosLibro/GustoMedio.png';
  import GustoAltoBajo from '../img/iconosLibro/GustoBajo.png';
  import CorazonVioleta from '../img/iconosLibro/CorazonVioletaCompleto.png'
  import CorazonAzul from '../img/iconosLibro/CorazonAzul.png'
  
  // Importación de estado global y librerías auxiliares
  import { useRutasGlobalRecursos } from '../stores/rutasGlobales';
  import { useCalificadosStore, useFavoritosStore, cargarNuevoDatoUsuarioXRecurso, datosUsuarioXrecurso } from '../stores/usuarioXrecurso'
  import { libreriaInstancia } from '../utils/axios';

  
  const userData = JSON.parse(localStorage.getItem('userDataEAT') || '{}');
  const IdUsuario : number = userData.Id || 0;
  const Calificados = useCalificadosStore();
  const Favoritos = useFavoritosStore();
  const listaCalificados = Calificados.Ids
  const listaFavoritos = Favoritos.Ids
  
  let nombreArchivo : string
  let esFavorito= false
  let direccionRec : string
  let IconoUrl = '';
  let corazon = ref(CorazonAzul);
  let calificacion  = ref(1)
  // Renderiza el libro al montar el componente
  onMounted(async () => {
    await nextTick(); // Espera a que la reactividad esté lista

    console.log(id);  // Verifica si `id` ahora está disponible
    const respuesta = await cargarRecurso(id);

    // Asignar el resultado a recurso.value, si no es undefined
    if (respuesta !== undefined) {
      console.log('Respuesta recibida:', respuesta);
      renderEpub(respuesta.Archivo);
      nombreArchivo = respuesta.Titulo
        if (respuesta.Calificacion >= 6.6) {
          IconoUrl = GustoAltoIcono; // Ruta importada
        }else if(respuesta.Calificacion >= 3.3){
          IconoUrl = GustoAltoMedio;
        }else{
          IconoUrl = GustoAltoBajo;
        }
      recurso.value = respuesta;  // Asigna el valor al ref usando .value
    }
  });
  const agregarFavorito = () =>{
    if (!esFavorito){
      corazon.value = CorazonVioleta
      esFavorito = true
      cargarNuevoDatoUsuarioXRecurso(IdUsuario,id, null, null,null)
      datosUsuarioXrecurso(IdUsuario)
    }else{
      corazon.value= CorazonAzul
      esFavorito = false
    }
    
    
  }
  const calificarRecurso = (nuevaCalificacion : number) =>{
    calificacion.value= nuevaCalificacion
    
  }
  function descargarLibro() {

    // Crear un enlace temporal para descargar el archivo
    const link = document.createElement('a');
    link.href = direccionRec;
    link.download = nombreArchivo; // Nombre sugerido para el archivo descargado
    document.body.appendChild(link);
    link.click(); // Simula el clic para iniciar la descarga
    document.body.removeChild(link); // Limpia el DOM eliminando el enlace
  }
  // Definición de propiedades
  const { id } = defineProps<{
    id: number;
  }>();
  if (listaFavoritos.some((num) => num === id)){
    
    corazon.value = CorazonVioleta
  }
  if (listaCalificados.some((num) => num === id)){
    const index = listaCalificados.findIndex((num) => num === id);
    const listaCalificaciones = Calificados.calificaciones
    calificacion.value = listaCalificaciones[index]
  }
  // Referencias y variables reactivas
  
  const recurso = ref<Recurso | undefined>(undefined);
  const containerRef = ref<HTMLDivElement | null>(null);
  const rutasStore = useRutasGlobalRecursos();
  const rutaBase = rutasStore.rutaRecursos;
  const toast = useToast();

  // Función para renderizar el libro en el contenedor
  const renderEpub = async (archivo: string) => {
    try {
      direccionRec = rutaBase + archivo;
      if (containerRef.value) {
        const book = ePub(direccionRec);
        const rendition = book.renderTo(containerRef.value, {
          width: '260px',
          height: '400px',
        });

        await book.ready;

        // Forzar estilos sin padding ni márgenes
        rendition.themes.default({
          '*': {
            "margin": '0 !important',
            "padding": '0 !important',
            "padding-left": '0 !important',
          },
        });
        rendition.themes.override('*', 'margin: 0 !important; padding: 0 !important;', true);
        
        rendition.display();
      }
    } catch (error) {
      console.error('Error al cargar el libro:', error);
    }
  };
  
  // Función para tomar el recurso desde el backend
  async function tomarRecurso(id: number): Promise<GraphQLResponse<{ tomarRecurso: Recurso }>> {
    console.log('Iniciar sesión ejecutado');
    try {
      const query = `
        query tomarRecurso($id: Int!) {
          tomarRecurso(id: $id) {
            Id
            Titulo
            Autor
            Categoria
            IdUsuario
            Formato
            Descripcion
            Archivo
            FechaOrigen
            Etiquetas
            Calificacion
            NumDescargas
          }
        }
      `;

      const variables = { id };

      const response = await libreriaInstancia.post('/query', { query, variables });
      return response.data;
    } catch (error) {
      toast.error('Error al tomar el recurso');
      console.error('Error al tomar el recurso:', error);
      throw error; // Lanzamos el error nuevamente
    }
  }
  // Función para cargar el recurso
  const cargarRecurso = async (id: number): Promise<Recurso | undefined> => {
    try {
      const respuesta = await tomarRecurso(id);
      if (respuesta.errors) throw new Error(respuesta.errors.map(e => e.message).join(', '));
      console.log(respuesta.data?.tomarRecurso);
      return respuesta.data?.tomarRecurso || undefined;
    } catch (error) {
      console.error('Error al obtener recursos:', error);
      toast.error('No se pudo conectar con el servidor.');
      return undefined;
    }
  };

</script>

<template>
  <div class="w-full h-full md:rounded-2xl bg-grisClaroEAT">
    <header class="flex justify-center items-center w-full h-14 md:justify-start md:pl-2 md:h-16">
      <BarraBusqueda nombreInput="" nombreSelec="" placeholder="Buscar. . ." :opciones="['Titulo','Autor','Etiquetas','Formato']" :valores="['titulo','autor','etiqueta','formato']"/>
    </header>
    <main class="w-full h-full scrollbar-hide ps-2 pe-2">
      <section class="flex flex-row w-4/6 h-5/6 rounded-xl shadow bg-grisEAT shadow-azulOscuroEAT">
        <article class="flex flex-col gap-10 items-center p-3 w-1/3">
          <div class="h-80" ref="containerRef"></div>
          <hr>
          <button @click="descargarLibro()" class="p-2 w-2/3 text-xl text-center text-white rounded-xl shadow cursor-pointer shadow-azulOscuroEAT bg-moradoOscuroEAT">{{  recurso?.Formato  }}</button>
          
        
        </article>
        <article class="flex flex-col flex-grow gap-5">
          <div class="flex justify-between p-2 w-full">
            <h1 class="text-2xl font-bold text-azulAunMasOscuroEAT">{{  recurso?.Titulo }}</h1>
            <img @click="agregarFavorito" :src="corazon" class="w-8 h-8 cursor-pointer" alt="">
          </div>
          <div class="flex gap-2">
            <h1 class="text-xl text-azulOscuroMedioEAT">Autor:</h1>
            <h1 class="text-xl text-azulGisaseoEAT">{{ recurso?.Autor }}</h1>
          </div>
          <div class="flex gap-2">
            <h1 class="text-xl">Clasificación:</h1>
            <h1 class="text-xl text-azulGisaseoEAT">{{ recurso?.Categoria }}</h1>
          </div class="">
          <div class="flex gap-2 items-center text-xl">
            <h1 class="text-xl">Etiquetas:</h1>
            <div class="flex flex-wrap gap-1">
              <span 
                v-for="(etiqueta, index) in recurso?.Etiquetas" 
                :key="index" 
                class="px-1 py-1 text-white rounded bg-azulAunMasOscuroEAT">
                {{ etiqueta }}
              </span>
            </div>
          </div>
          <p class="text-lg text-azulGisaseoEAT">{{ recurso?.Descripcion }}</p>
          <div class="flex gap-2 items-center">
            <h1 class="text-xl">Calificacion:</h1>
            <img class="mt-1 w-8 h-8 rounded" :src=IconoUrl alt="">
          </div>
          <div class="flex gap-2 items-center">
            <h1 class="text-xl">Calificar:</h1>
            <img 
            @click="calificarRecurso(10)"
            :class = "{'shadow-moradoOscuroEAT shadow': calificacion === 10}"
            class="mt-1 w-8 h-8 rounded-full cursor-pointer active:shadow-md active:shadow-moradoOscuroEAT" :src="GustoAltoIcono" alt="">
            <img 
            @click="calificarRecurso(5)"
            :class = "{'shadow-moradoOscuroEAT shadow': calificacion === 5}"
            class="mt-1 w-8 h-8 rounded-full cursor-pointer active:shadow-md active:shadow-moradoOscuroEAT" :src="GustoAltoMedio" alt="">
            <img 
            @click="calificarRecurso(0)"
            :class = "{'shadow-moradoOscuroEAT shadow': calificacion === 0}"
            class="mt-1 w-8 h-8 rounded-full cursor-pointer active:shadow-md active:shadow-moradoOscuroEAT" :src="GustoAltoBajo" alt="">
          </div>
        </article>
      </section>
      <section>
        
      </section>
    </main>
  </div>
</template>