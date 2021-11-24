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

## Frontend

Minimal page built using React. Project was created with command `yarn create vite docker-workshop --template react-ts`

```sh
# Build image so it can be run as container
cd frontend
docker build -t docker-workshop-frontend .
```

```sh
# Run frontend container
# -p <host>:<container> stands for port forwarding
# it makes port 80 inside container available on host port 3000
# so you can open http://localhost:3000 and see the application
docker run --rm -p 3000:80 docker-workshop-frontend
```
