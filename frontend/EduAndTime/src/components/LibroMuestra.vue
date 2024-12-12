<script setup lang="ts">
import { ref, onMounted } from 'vue';
import ePub from 'epubjs';
import GustoAltoIcono from '../img/iconosLibro/GustoAlto.png';
import GustoAltoMedio from '../img/iconosLibro/GustoMedio.png';
import GustoAltoBajo from '../img/iconosLibro/GustoBajo.png';





// Referencia para el contenedor donde se renderizará el libro
const containerRef = ref<HTMLDivElement | null>(null);

// Desestructura las propiedades
const {id, titulo, autor, categoria, formato, epubUrl, etiquetas, calificacion } = defineProps<{
  id: number;
  titulo: string;
  autor: string;
  categoria: string;
  formato: string;
  epubUrl: string;
  etiquetas: string[]; // Array de strings
  calificacion: number; // Número (entero o flotante)
}>();

let IconoUrl = '';
if (calificacion >= 6.6) {
  IconoUrl = GustoAltoIcono; // Ruta importada
}else if(calificacion >= 3.3){
  IconoUrl = GustoAltoMedio;
}else{
  IconoUrl = GustoAltoBajo;
}



const renderEpub = async () => {
  try {
    if (containerRef.value) {
      const book = ePub(epubUrl);
      const rendition = book.renderTo(containerRef.value, {
        width: '120px',
        height: '200px',
      });

      await book.ready;

      // Forzar estilos sin padding ni márgenes
      rendition.themes.default({
        '*': {
          "margin": '0 !important',
          "padding": '0 !important',
          "padding-left": '0 !important'
        },
      });
      rendition.themes.override('*', 'margin: 0 !important; padding: 0 !important;', true);
      

      rendition.display();
    }
  } catch (error) {
    console.error('Error al cargar el libro:', error);
  }
  
};


const emit = defineEmits(['click']);

const handleClick = () => {
  emit('click', id); // Emitimos el id al componente padre
};
onMounted(renderEpub);
</script>

<template>
  <div class="flex pt-2.5 w-11/12 h-auto bg-white rounded-md border-2 cursor-pointer border-azulAunMasOscuroEAT md:w-auto" @click="handleClick()">
    <div class="h-40" ref="containerRef"></div>
    <div class="flex flex-col justify-center w-full">
      <h2 class="font-bold text-azulAunMasOscuroEAT text-md">{{ titulo }}</h2>
      <h3 class="text-azulGisaseoEAT">{{ autor }}</h3>
      <div class="flex">
        <h3 class="flex-grow text-azulAunMasOscuroEAT">{{ categoria }}</h3>
        <h3 class="mr-2 text-azulGisaseoEAT">{{ formato }}</h3>
      </div>

      <!-- Renderizar etiquetas con v-for -->
      <div class="flex flex-wrap gap-1">
        <span 
          v-for="(etiqueta, index) in etiquetas" 
          :key="index" 
          class="px-1 py-1 text-white rounded bg-azulAunMasOscuroEAT">
          {{ etiqueta }}
        </span>
          
      </div>
      <img class="mt-1 w-7 h-7" :src=IconoUrl alt="">
      
    </div>
  </div>
</template>

