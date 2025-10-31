import { http, API } from './http'

export async function login({ email, password }) {
    const res = await http.post(`${API.USERS}/auth/login`, { email, password })
    return res.data
}
