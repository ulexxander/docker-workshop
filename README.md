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

## MongoDB

When we have `docker-compose` adding new containers is very easy.

Few specifics of database containers:

1. We probably don't want our containers to be **exposed** to the rest world - we don't forward ports and we don't need to configure reverse-proxy for them.
2. We may still use **ports forwarding** during development (database runs inside docker container, but our app usually runs on localhost)
3. Database containers save state to filesystem - to persist that data reliably we should use **Docker volumes**.
4. We want our backend apps to run on same **Docker network** with database containers, so we can connect them using **hostnames** and **exposed** (not forwarded) ports.

## Deploying to production

Because we have all our infrastructure defined in `docker-compose.yml` file - it is all we need to run our stack!

Except images - some of them are public (traefik, mongo) and some of them are our custom images (docker-workshop-backend, docker-workshop-frontend).

Public images are pulled from [Docker Hub](https://hub.docker.com/).
Our images were available only locally (we produced them with `docker build`).

So to **pull** our custom images in production environment we first need to **push** them to some registry (just like with git)

We have more options:

- [Docker Hub](https://hub.docker.com/) - default and most popular registry (best for open-source projects, but also supports private repositories)
- [GitLab Container Registry](https://docs.gitlab.com/ee/user/packages/container_registry) - my choise for non-open-source projects (images are protected and can be pulled only with access tokens)
- [GitHub Container registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry) - similar thing

We will store images in **Docker Hub** in this workshop.

```sh
# Build images that will be published
cd backend
docker build -t ulexxander/docker-workshop-backend:1.0.0 .
cd ../frontend
docker build -t ulexxander/docker-workshop-frontend:1.0.0 .
```

```sh
# Authenticate to the images registry, you will be prompted for username and password
docker login
```

```sh
# Push images to the registry
docker push ulexxander/docker-workshop-backend:1.0.0
docker push ulexxander/docker-workshop-frontend:1.0.0
```

```sh
# Now we can pull images anywhere registry is accessible!
# Like on our production environment!
docker pull ulexxander/docker-workshop-backend:1.0.0
docker pull ulexxander/docker-workshop-frontend:1.0.0
```

## Workshop Checkpoints

You can follow progress with `git checkout <hash>`

| Commit Message                                                          | Hash                                     |
| ----------------------------------------------------------------------- | ---------------------------------------- |
| prepared Dockerfile for backend and explained docker commands in README | 71fe95659f05f057bd6a26bc24c3345481fd8c6b |
| frontend Dockerfile and more docker commands in README                  | 4d488562d2b0d726ad5f00bcc8eaff7602c5c129 |
| simple docker-compose to run both backend and frontend                  | 709d8ab1e40f7a7f85c76fbaf16b7d7526f707ca |
| readme - added Traefik and docker-compose section                       | 2af80721db0122b046bbd4bf3da3cdcf67cfa0a3 |
| backend - implemented NotesStore backed by mongodb and tested it        | 41b520aaf297ee9a90e77ab8fea81575c6029b39 |
| docker-compose mongo service - added volume to persist database         | 7a974019c97bc8f353ead1f0a39b40d9da1eab98 |
| readme - fixed Now we can pull images section                           | 5ad5b74e22f373b676efe84f280d67621845b500 |
