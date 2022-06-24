FROM node:lts-bullseye AS build-step

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .
RUN npm run build

# production environment
FROM nginx:mainline

COPY ./nginx/nginx.conf /etc/nginx/nginx.conf
COPY --from=build-step /app/build /frontend/build

# EXPOSE 80
# CMD ["nginx", "-g", "daemon off;"]