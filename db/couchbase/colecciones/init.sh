#!/bin/bash
sleep 10
COUCHBASE_ADMIN_USERNAME=${COUCHBASE_ADMIN_USERNAME:-"administrador"}
COUCHBASE_ADMIN_PASSWORD=${COUCHBASE_ADMIN_PASSWORD:-"tearsthemoon"}
COUCHBASE_BUCKET=${COUCHBASE_BUCKET:-"EduAndTime"}
COUCHBASE_RAM_QUOTA=${COUCHBASE_RAM_QUOTA:-"512"}
COUCHBASE_INDEX_RAM_QUOTA=${COUCHBASE_INDEX_RAM_QUOTA:-"256"}
COUCHBASE_FTS_RAM_QUOTA=${COUCHBASE_FTS_RAM_QUOTA:-"512"}
COUCHBASE_EVENTING_RAM_QUOTA=${COUCHBASE_EVENTING_RAM_QUOTA:-"256"}
COUCHBASE_QUERY_RAM_QUOTA=${COUCHBASE_QUERY_RAM_QUOTA:-"512"}
COUCHBASE_HOST=${COUCHBASE_HOST:-"127.0.0.1"}

# Esperar a que Couchbase Server esté disponible
until curl -s http://$COUCHBASE_HOST:8091/pools >/dev/null; do
  echo "Esperando a que Couchbase Server esté disponible..."
  sleep 5
done

echo "Couchbase Server está listo, iniciando configuración inicial..."

# Inicializar el clúster de Couchbase
couchbase-cli cluster-init \
  --cluster 127.0.0.1:8091 \
  --cluster-username $COUCHBASE_ADMIN_USERNAME \
  --cluster-password $COUCHBASE_ADMIN_PASSWORD \
  --cluster-ramsize $COUCHBASE_RAM_QUOTA \
  --cluster-index-ramsize $COUCHBASE_INDEX_RAM_QUOTA \
  --cluster-fts-ramsize $COUCHBASE_FTS_RAM_QUOTA \
  --cluster-eventing-ramsize $COUCHBASE_EVENTING_RAM_QUOTA \
  --cluster-query-ramsize $COUCHBASE_QUERY_RAM_QUOTA \
  --services data,index,query,fts,eventing

# Crear un bucket
curl -v -X POST http://$COUCHBASE_HOST:8091/pools/default/buckets \
  -u $COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD \
  -d name=$COUCHBASE_BUCKET \
  -d ramQuota=512 \
  -d authType=sasl \
  -d saslPassword=$COUCHBASE_ADMIN_PASSWORD

# Opcional: Crear un usuario para acceder al bucket
curl -v -u $COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD \
     -X PUT http://$COUCHBASE_HOST:8091/settings/rbac/users/local/$COUCHBASE_ADMIN_USERNAME \
     -d password=$COUCHBASE_ADMIN_PASSWORD \
     -d roles=bucket_full_access[$COUCHBASE_BUCKET],admin

echo "Couchbase ha sido configurado exitosamente."

# Crear un scope dentro del bucket
curl -v -u $COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD \
     -X POST http://$COUCHBASE_HOST:8091/pools/default/buckets/$COUCHBASE_BUCKET/scopes \
     -d name=usuarios

curl -v -u $COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD \
     -X POST http://$COUCHBASE_HOST:8091/pools/default/buckets/$COUCHBASE_BUCKET/scopes \
     -d name=documentos

curl -v -u $COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD \
     -X POST http://$COUCHBASE_HOST:8091/pools/default/buckets/$COUCHBASE_BUCKET/scopes \
     -d name=chats     

# Crear una colección dentro del scope
curl -v -u $COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD \
     -X POST http://$COUCHBASE_HOST:8091/pools/default/buckets/$COUCHBASE_BUCKET/scopes/usuarios/collections \
     -d name=usuario

curl -v -u $COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD \
     -X POST http://$COUCHBASE_HOST:8091/pools/default/buckets/$COUCHBASE_BUCKET/scopes/usuarios/collections \
     -d name=configProductividad

curl -v -u $COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD \
     -X POST http://$COUCHBASE_HOST:8091/pools/default/buckets/$COUCHBASE_BUCKET/scopes/chats/collections \
     -d name=chatGrupo

curl -v -u $COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD \
     -X POST http://$COUCHBASE_HOST:8091/pools/default/buckets/$COUCHBASE_BUCKET/scopes/chats/collections \
     -d name=chatMensajes
  
curl -v -u $COUCHBASE_ADMIN_USERNAME:$COUCHBASE_ADMIN_PASSWORD \
    -X POST http://$COUCHBASE_HOST:8091/pools/default/buckets/$COUCHBASE_BUCKET/scopes/documentos/collections \
    -d name=documento   

echo "Couchbase ha sido configurado con buckets, scopes, colecciones."



