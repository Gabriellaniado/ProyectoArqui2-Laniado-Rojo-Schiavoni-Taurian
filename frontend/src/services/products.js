import { http, API } from './http'

// Lista productos desde el Search/List API (Solr)
export async function listProducts({ q } = {}) {
  const res = await http.get(`${API.SEARCH}/items`, { params: { q } })
  const data = res.data
  const items = Array.isArray(data?.items)
    ? data.items
    : (Array.isArray(data?.results) ? data.results : [])
  return { items }
}

// Obtiene detalle desde Products API (Mongo)
export async function getProductById(id) {
  const res = await http.get(`${API.PRODUCTS}/items/${id}`)
  return res.data?.item ?? res.data
}
