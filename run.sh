#replaces the code and reruns the docker
docker-compose up &
echo up
docker exec -it http-docker-app-1 go build -v -o /usr/local/bin/app ./...
echo build completed
docker exec -it http-docker-app-1 kill 1