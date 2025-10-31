# Build stage
FROM node:20-alpine AS build
WORKDIR /app

# Instalar deps primero para cache eficiente
COPY package*.json ./
RUN npm ci --no-audit --no-fund

# Copiar el resto del código
COPY . .

# Permitir inyectar variables VITE_ al momento del build
ARG VITE_APP_VERSION
ARG VITE_USERS_API
ARG VITE_SEARCH_API
ARG VITE_PRODUCTS_API
ENV VITE_APP_VERSION=${VITE_APP_VERSION}
ENV VITE_USERS_API=${VITE_USERS_API}
ENV VITE_SEARCH_API=${VITE_SEARCH_API}
ENV VITE_PRODUCTS_API=${VITE_PRODUCTS_API}

# Build de producción
RUN npm run build

# Runtime stage: Nginx estático con SPA fallback
FROM nginx:1.27-alpine AS runtime
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
HEALTHCHECK CMD wget -qO- http://localhost/ || exit 1

