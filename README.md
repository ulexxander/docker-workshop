# Docker Workshop

Practical example of using Docker for those who learn it.

## Backend

Simple serice written in Go.

```sh
# Build image so it can be run as container
cd backend
docker build -t docker-workshop-backend .
```

```sh
# Run image (starts new container)
# --rm flag can be used to delete container after it stops
docker run docker-workshop-backend
```

```sh
# List all containers (-a means running AND stopped)
docker ps -a
```

```sh
# Remove container
docker container rm <container_id>
```
