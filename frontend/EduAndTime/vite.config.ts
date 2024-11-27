import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import vueJsx from '@vitejs/plugin-vue-jsx';

export default defineConfig({
  plugins: [vue(), vueJsx()],
  build: {
    sourcemap: true, // Habilita los mapas fuente
  },
  server: {
    sourcemapIgnoreList: false, // Opcional: Asegúrate de que no se ignoren
  },
});
