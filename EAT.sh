#!/bin/dash

# Configurar nombres para reutilización
COMPOSE_FILE="docker-compose.yml"
CONTAINER_NAME="couchbase"
INIT_SCRIPT="/init.sh"
INDEX_SCRIPT="/index.sh"

# Paso 1: Reconstruir las imágenes sin caché
echo "Reconstruyendo imágenes sin caché..."
if ! docker-compose -f "$COMPOSE_FILE" build --no-cache; then
    echo "Error al construir las imágenes. Saliendo..."
    exit 1
fi

# Paso 2: Iniciar los servicios
echo "Iniciando servicios..."
if ! docker-compose up -d; then
    echo "Error al iniciar los servicios. Saliendo..."
    exit 1
fi

# Paso 3: Esperar que Couchbase esté listo
echo "Esperando a que Couchbase esté listo..."
sleep 26

# Paso 4: Ejecutar el script de inicialización
echo "Ejecutando script de inicialización en Couchbase..."
if ! docker exec -it "$CONTAINER_NAME" "$INIT_SCRIPT"; then
    echo "Error al ejecutar $INIT_SCRIPT en $CONTAINER_NAME. Saliendo..."
    exit 1
fi

# Paso 5: Esperar un poco antes de ejecutar el script de índices
echo "Esperando antes de ejecutar el script de índices..."
sleep 30

# Paso 6: Ejecutar el script de índices
echo "Ejecutando script de índices en Couchbase..."
if ! docker exec -it "$CONTAINER_NAME" "$INDEX_SCRIPT"; then
    echo "Error al ejecutar $INDEX_SCRIPT en $CONTAINER_NAME. Saliendo..."
    exit 1
fi

echo "Scripts ejecutados exitosamente."
exit 0
