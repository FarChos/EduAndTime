/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}'
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
        azulOscuroEAT: '#1c2833',
        azulClaroEAT: '#11538a',
        verdeOscuroEAT: '#0b5345',
        moradoClaro: '#d600f0',
        moradoOscuroEAT: '#8e44ad',
        verdeClaroEAT: '#2ecc71',
        grisEAT: '#aeb6bf'
      },
    },
  },
  plugins: [],
}
