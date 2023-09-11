docker restart Etcd-server
docker-compose build
docker-compose up -d
docker network connect apim-with-analytics_default login
docker network connect apim-with-analytics_default Etcd-server
