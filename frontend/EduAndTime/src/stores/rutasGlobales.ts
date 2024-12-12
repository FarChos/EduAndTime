import { defineStore } from 'pinia'

export const useRutasGlobalImg = defineStore('rutaImgPerf', {
  state: () => ({
    rutaImagenesPerfil: '/public/imagenesPerfil/'
  }),
})

export const useRutasGlobalRecursos = defineStore('rutaRecursos', {
  state: () => ({
    rutaRecursos: '/public/recursos/'
  }),
})
