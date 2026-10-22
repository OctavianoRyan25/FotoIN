# Multi-stage Dockerfile for FotoIN Dashboard
FROM node:20-alpine AS builder

WORKDIR /app

# Cache package dependencies
COPY package*.json ./
RUN npm ci

# Copy sources and build
COPY . .
RUN npm run build

# Nginx runtime stage
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
