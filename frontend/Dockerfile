FROM node:14 AS builder
WORKDIR /build
COPY package.json yarn.lock ./
RUN yarn
COPY . .
RUN yarn build

FROM nginx
WORKDIR /usr/share/nginx/html
COPY --from=builder /build/dist .
RUN \
  sed -i "s/location \/ {$/location \/ {\n try_files \$uri \/index.html;/g" /etc/nginx/conf.d/default.conf
CMD ["nginx", "-g", "daemon off;"]
