// src/router.ts
import { createRouter, createWebHistory } from 'vue-router'
import sesion from './layouts/iniciosesion.vue';
import registro from './layouts/registro.vue';
import home from './layouts/home.vue';
const routes = [
  {
    path: '/',
    name: 'sesion',
    component: sesion
  },
  {
    path: '/registro',
    name: 'registro',
    component: registro
  },
  {
    path: '/home',
    name: 'home',
    component: home
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
