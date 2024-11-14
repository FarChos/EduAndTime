# Etapa 1: Construcción de la aplicación
FROM opensuse/tumbleweed:latest AS builder

# Instalar dependencias necesarias
RUN zypper refresh && zypper install -y go git

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar los archivos go.mod y go.sum
COPY go.mod go.sum ./

# Descargar las dependencias
RUN go mod download

# Copiar el código fuente
COPY . .

# Establecer la arquitectura
ENV GOOS=linux
ENV GOARCH=amd64

# Compilar la aplicación
RUN go build -o EduAndTime .

# Cambiar permisos para el binario
RUN chmod +x /app/EduAndTime

# Etapa 2: Crear la imagen final
FROM opensuse/tumbleweed:latest

# Instalar certificados para HTTPS
RUN zypper refresh && zypper install -y ca-certificates file

# Copiar el binario desde la etapa de construcción
COPY --from=builder /app/EduAndTime /EduAndTime

# Dar permisos de ejecución al binario en la etapa final
RUN chmod +x /EduAndTime



# Establecer el comando por defecto
CMD ["./EduAndTime"]

# Exponer el puerto que usa la aplicación
#EXPOSE 8080

