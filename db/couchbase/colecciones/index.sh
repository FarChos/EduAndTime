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
    curl -u "$COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD" -X POST http://localhost:8093/query/service \
    -H "Content-Type: application/json" \
    -d '{"statement": "SELECT * FROM `eduandtime` LIMIT 1;"}'
    if [ $? -eq 0 ]; then
        break
    fi
    sleep 5
done

echo "Couchbase Server y N1QL están listos, iniciando creación de los índices..."

# Crear índices
curl -u "$COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD" \
     -X POST http://localhost:8093/query/service \
     -d "statement=CREATE INDEX \`idx_idDocumento\` ON \`$COUCHBASE_BUCKET\`.\`documentos\`.\`documento\`(idDocumento)"

curl -u "$COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD" \
     -X POST http://127.0.0.1:8093/query/service \
     -d "statement=CREATE INDEX \`idx_idChatGrupo\` ON \`$COUCHBASE_BUCKET\`.\`chats\`.\`chatGrupo\`(idChat)"

curl -u "$COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD" \
     -X POST http://127.0.0.1:8093/query/service \
     -d "statement=CREATE INDEX \`idx_idChatMensajes\` ON \`$COUCHBASE_BUCKET\`.\`chats\`.\`chatMensajes\`(idChat)"

curl -u "$COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD" \
     -X POST http://127.0.0.1:8093/query/service \
     -d "statement=CREATE INDEX \`idx_idUsuario\` ON \`$COUCHBASE_BUCKET\`.\`usuarios\`.\`usuario\`(idUsuario)"

curl -u "$COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD" \
     -X POST http://127.0.0.1:8093/query/service \
     -d "statement=CREATE INDEX \`idx_idUsuarioConfig\` ON \`$COUCHBASE_BUCKET\`.\`usuarios\`.\`configProductividad\`(idUsuario)"  

echo "Los índices han sido configurados :)"
