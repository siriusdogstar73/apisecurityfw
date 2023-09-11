docker restart Etcd-server
docker-compose build
docker-compose up -d
docker network connect apim-with-analytics_default error
docker network connect apim-with-analytics_default Etcd-server
