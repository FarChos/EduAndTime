// src/router.ts
import { createRouter, createWebHistory, RouteLocationNormalized } from 'vue-router';
import sesion from './layouts/iniciosesion.vue';
import registro from './layouts/registro.vue';
import menu from './layouts/menu.vue';
import welcome from './layouts/welcome.vue';
import libreriaPrincipal from './layouts/libreriaPrincipal.vue';
import libroCompleto from './layouts/libroCompleto.vue';
import subirRecurso from './layouts/subirRecurso.vue';
import cuenta from './layouts/cuenta.vue';
import favoritos from './layouts/favoritos.vue';
import misRecursos from './layouts/misRecursos.vue';
import editarPerfil from './layouts/editarPrf.vue'
import eliminarPerfil from './layouts/eliminarPrf.vue'

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
    redirect: '/menu/welcome',
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
        path: 'subirRecurso',
        name: 'subirRecurso',
        component: subirRecurso,
      },
      {
        path: 'libroCompleto',
        name: 'libroCompleto',
        component: libroCompleto,
        props: (route: RouteLocationNormalized) => ({ id: Number(route.query.id) }),
      },
    ],
  },
  {
    path: '/cuenta',
    name: 'cuenta',
    component: cuenta,
    //redirect: '/menu/welcome',
    children: [
      {
        path: 'favoritos', // Subruta de "/menu"
        name: 'favoritos',
        component: favoritos,
      },
      {
        path: 'misRecursos', // Subruta de "/menu"
        name: 'misRecursos',
        component: misRecursos,
      },
      {
        path: 'editarPerfil', // Subruta de "/menu"
        name: 'editarPerfil',
        component: editarPerfil,
      },
      {
        path: 'eliminarPerfil', // Subruta de "/menu"
        name: 'eliminarPerfil',
        component: eliminarPerfil,
      },
    ]
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
