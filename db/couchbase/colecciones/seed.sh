#!/bin/bash

COUCHBASE_ADMIN_USERNAME="${COUCHBASE_ADMIN_USERNAME:-"administrador"}"
COUCHBASE_ADMIN_PASSWORD="${COUCHBASE_ADMIN_PASSWORD:-"tearsthemoon"}"
COUCHBASE_BUCKET="${COUCHBASE_BUCKET:-"EduAndTime"}"

# Esperar hasta que Couchbase esté disponible
until curl -s http://localhost:8091/pools >/dev/null; do
  echo "Esperando a que Couchbase Server esté disponible..."
  sleep 5
done

# Esperar hasta que N1QL (puerto 8093) esté disponible
echo "Esperando a que el servicio N1QL esté disponible en el puerto 8093..."
for i in {1..10}; do
    response=$(curl -u "$COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD" -X POST http://localhost:8093/query/service \
    -H "Content-Type: application/json" \
    -d '{"statement": "SELECT * FROM `'$COUCHBASE_BUCKET'` LIMIT 1;"}')
    
    if echo "$response" | grep -q '"results":\s*\[\]'; then
        echo "N1QL está disponible."
        break
    fi
    echo "Esperando que N1QL esté listo..."
    sleep 5
done

if [ $? -ne 0 ]; then
    echo "N1QL no está disponible después de 10 intentos."
    exit 1
fi

echo "Couchbase Server y N1QL están listos, iniciando creación de los datos semilla..."

curl -u "$COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD" \
     -X POST http://localhost:8093/query/service \
     -d 'statement=INSERT INTO `'$COUCHBASE_BUCKET'`.`usuario` (KEY, VALUE) VALUES ("usuario::1", { "tipo": "usuario", "estaActivo": true, "version": 1, "idUsuario": 1, "docCalificados": [{"idDoc":1,"calificacion":4}], "docOriginados": [{"idDoc":1}]})'

exit 0