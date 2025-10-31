# Frontend (React + Vite)

App React para el e-commerce de mates.

## Requisitos
- Node 18+ y npm
- Backends levantados (ver puertos en `.env.development`)

## Variables de entorno
Configurar los endpoints en `.env.development` (copiar de `.env.example` si hace falta):

```
VITE_USERS_API=http://localhost:8080
VITE_SEARCH_API=http://localhost:8081
VITE_PRODUCTS_API=http://localhost:8082
```

## Comandos
- `npm install`
- `npm run dev` (dev server)
- `npm run build` (build producción)
- `npm run preview` (previsualizar build local)

## Notas
- Listado de productos usa Search/List API (`/items`).
- Detalle de producto usa Products API (`/items/:id`).

## Docker

- Build y run con Docker Compose (puerto 8088):
  - `docker compose -f docker-compose.yml build`
  - `docker compose -f docker-compose.yml up -d`
  - App: `http://localhost:8088`

- Endpoints configurables en tiempo de build (usando args VITE_):
  - `VITE_USERS_API` (default `http://host.docker.internal:8080`)
  - `VITE_SEARCH_API` (default `http://host.docker.internal:8081`)
  - `VITE_PRODUCTS_API` (default `http://host.docker.internal:8082`)

Si tus APIs corren en el host, `host.docker.internal` permite acceder desde el contenedor a servicios del host en Windows/macOS. En Linux podés reemplazarlo por la IP del host.
