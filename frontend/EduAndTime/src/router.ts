// src/router.ts
import { createRouter, createWebHistory } from 'vue-router';
import sesion from './layouts/iniciosesion.vue';
import registro from './layouts/registro.vue';
import menu from './layouts/menu.vue';
import welcome from './layouts/welcome.vue';
import libreriaPrincipal from './layouts/libreriaPrincipal.vue';
import libroCompleto from './layouts/libroCompleto.vue';

const routes = [
  {
    path: '/',
    name: 'sesion',
    component: sesion,
  },
  {
    path: '/registro',
    name: 'registro',
    component: registro,
  },
  {
    path: '/menu',
    name: 'menu',
    component: menu,
    //redirect: '/menu/welcome',
    children: [
      {
        path: 'welcome', // Subruta de "/menu"
        name: 'welcome',
        component: welcome,
      },
      {
        path: 'libreriaPrincipal',
        name: 'libreriaPrincipal',
        component: libreriaPrincipal,
      },
      {
        path: 'libroCompleto',
        name: 'libroCompleto',
        component: libroCompleto,
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
