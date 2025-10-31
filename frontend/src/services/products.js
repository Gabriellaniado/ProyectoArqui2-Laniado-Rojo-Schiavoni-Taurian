import { http, API } from './http'

export async function listProducts({ q } = {}) {
    const res = await http.get(`${API.PRODUCTS}/products`, { params: { q } })
    return res.data // esperado: { items: [...] }
}

export async function getProductById(id) {
    const res = await http.get(`${API.PRODUCTS}/products/${id}`)
    return res.data
}
