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

RUN apk add --no-cache openssl

RUN rm -rf /usr/share/nginx/html/*

COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf.template
RUN chmod +x /docker-entrypoint.sh

ENV BACKEND_HOST=api
ENV LIVEKIT_HOST=livekit