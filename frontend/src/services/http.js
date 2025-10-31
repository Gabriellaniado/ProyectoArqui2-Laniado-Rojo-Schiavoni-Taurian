import axios from 'axios'

// Configurá tus MS por .env
export const API = {
    USERS: import.meta.env.VITE_USERS_API,       // ej: http://localhost:8081
    PRODUCTS: import.meta.env.VITE_PRODUCTS_API, // ej: http://localhost:8082
}

export const http = axios.create({ timeout: 8000 })

http.interceptors.request.use(cfg => {
    // si usás cookie HttpOnly para JWT, no metas el token aquí.
    cfg.headers['X-App-Version'] = import.meta.env.VITE_APP_VERSION ?? 'dev'
    return cfg
})
