#/bin/sh

export PORT=50051
export URL=https://youtube.com
export SERVER_ADDRESS="0.0.0.0:8080"

envsubst < ./proxy/nginx.template > ./proxy/nginx.conf
envsubst < ./docker-compose.template > ./docker-compose.yaml
docker-compose up -d --build