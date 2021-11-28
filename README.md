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

## Traefik and docker-compose

> Traefik is an open-source Edge Router that makes publishing your services a fun and easy experience. It receives requests on behalf of your system and finds out which components are responsible for handling them.

> Compose is a tool for defining and running multi-container Docker applications. With Compose, you use a YAML file to configure your application’s services. Then, with a single command, you create and start all the services from your configuration.

`docker-compose.yml` defines 3 services:

- `traefik` - accepts all http requests and forwards them based on domain, it also has its admin dashboard
- `backend` - our backend image that has port 80 exposed (see Dockerfile) and labels defined so traefik can proxy it
- `frontend` - our frontend image that calls `/api` endpoints which are fulfilled by backend

```sh
# Run whole stack (make sure you have images already built)
docker-compose up -d
```

Open traefik dashboard: http://traefik.localhost

Backend is exposed at standalone domain: http://backend.docker-workshop.localhost or at http://frontend.docker-workshop.localhost/api

Frontend is served at http://frontend.docker-workshop.localhost
