# Etapa 1: Construcción
FROM node:latest AS build

# Establece el directorio de trabajo
WORKDIR /app

# Copia los archivos de configuración y el código fuente
COPY package*.json ./
COPY . .

# Instala las dependencias
RUN npm install

# Construye el proyecto
RUN npm run build

# Etapa 2: Servir
#FROM nginx:latest

# Copiar la configuración personalizada de Nginx
#COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copiar archivos de construcción al directorio de Nginx
#COPY --from=build /app/dist /usr/share/nginx/html

# Exponer el puerto que usará Nginx
#EXPOSE 80
EXPOSE 5173

# Comando por defecto
#CMD ["nginx", "-g", "daemon off;"]
CMD ["npm", "run", "dev"]
