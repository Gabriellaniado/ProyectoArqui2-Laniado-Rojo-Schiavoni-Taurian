import axios from 'axios'

// Configurá tus MS por .env
export const API = {
  // Users (login/registro)
  USERS: import.meta.env.VITE_USERS_API, // ej: http://localhost:8080
  // Products API (detalle / CRUD items)
  PRODUCTS: import.meta.env.VITE_PRODUCTS_API, // ej: http://localhost:8082
  // Search/List API (búsqueda y listado paginado)
  SEARCH: import.meta.env.VITE_SEARCH_API ?? import.meta.env.VITE_PRODUCTS_API,
}

export const http = axios.create({ timeout: 8000 })

http.interceptors.request.use(cfg => {
  // si usás cookie HttpOnly para JWT, no metas el token acá.
  cfg.headers['X-App-Version'] = import.meta.env.VITE_APP_VERSION ?? 'dev'
  return cfg
})

