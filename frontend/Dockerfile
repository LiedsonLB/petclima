# ---------- BUILD ----------
FROM node:22-alpine AS build

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

ARG VITE_API_URL
ARG VITE_APP_KEY
ARG VITE_LIVEKIT_URL

ENV VITE_API_URL=$VITE_API_URL
ENV VITE_APP_KEY=$VITE_APP_KEY
ENV VITE_LIVEKIT_URL=$VITE_LIVEKIT_URL

RUN npm run build

# ---------- PRODUCTION ----------
FROM nginx:alpine

RUN rm -rf /usr/share/nginx/html/*

COPY --from=build /app/dist /usr/share/nginx/html
# A imagem oficial do nginx roda envsubst em /etc/nginx/templates/*.template
# e grava o resultado em /etc/nginx/conf.d/ (por isso o arquivo é só o bloco server).
COPY nginx.conf /etc/nginx/templates/default.conf.template

ENV BACKEND_HOST=api
ENV LIVEKIT_HOST=livekit