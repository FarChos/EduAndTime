/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{html,js,vue,jsx,tsx}'
  ],
  
  theme: {
    extend: {
      fontFamily: {
        gabriela: ['Gabriela', 'serif'],
        abhayaRegular: ['AbhayaLibre-Regular', 'serif'],
        abhayaBold: ['AbhayaLibre-Bold', 'serif'],
      },
      colors: {
        transparent: 'transparent',
        azulProfundo: '#1c2830',
        azulOscuroEAT: '#2C3E50',
        azulOscuroMedioEAT: '#1c2832',
        azulAunMasOscuroEAT: '#24313E',
        azulClaroEAT: '#11538a',
        azulGisaseoEAT: '#85929E',
        verdeOscuroEAT: '#0b5345',
        moradoClaro: '#d600f0',
        moradoOscuroEAT: '#8e44ad',
        verdeClaroEAT: '#2ecc71',
        grisEAT: '#D9E1E3',
        grisClaroEAT: '#ECF0F1'
      },
    },
  },
  plugins: [],
}
