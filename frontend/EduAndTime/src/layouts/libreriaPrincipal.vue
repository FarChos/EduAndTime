<script setup lang="ts">
import { ref } from 'vue';
import LibroMuestra from '../components/LibroMuestra.vue';
import BotonFiltroCategoria from '../components/BotonFiltroCategoria.vue';
import BarraBusqueda from '../components/BarraBusqueda.vue';


// Estado inicial del botón enfocado
const focusedButton = ref<string | null>(null);
const inicial = 'general'
focusedButton.value = inicial
// Manejar el cambio de foco en los botones
function cambiarFocus(button: string) {
  focusedButton.value = button;
}
</script>

<template>
  <div class="flex flex-col w-full h-full bg-grisClaroEAT md:rounded-2xl">
    <header class="flex justify-center items-center w-full h-14 md:justify-start md:pl-2 md:h-16">
      <img class="absolute left-5 z-10 h-6 md:left-24" src="../img/libreria/iconoLupa.png" alt="">
      <!-- Input -->
      <BarraBusqueda nombreInput="" nombreSelec="" placeholder="Buscar. . ." :opciones="['Titulo','Autor','Etiquetas','Formato']" :valores="['','','','']"/>

    </header>
    <hr class="border-azulOscuroEAT">
    <nav class="flex overflow-x-auto gap-1 items-center pl-1 h-12 whitespace-nowrap scrollbar-hide md:justify-center">
      <BotonFiltroCategoria label="+" nombre="+" :botonFocus="focusedButton" @changeFocus="cambiarFocus"/>
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
        <LibroMuestra
          epubUrl="/public/recursos/Los novios del invierno.epub"
          titulo="Los Jardines de la Luna"
          autor ="Steven Erikson"
          categoria= "Arte"
          formato="EPUB"
          :etiquetas="['Fantasía', 'Épico', 'Magia']"
          :calificacion= 9
        />
      </section>
      <h1 class="md:text-xl text-azulOscuroEAT">Recomendados</h1>
      <section>
        
      </section>
    </main>
  </div>

</template>