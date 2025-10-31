import { useState } from 'react'
import { login } from '../services/users'

export default function Login() {
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [msg, setMsg] = useState('')

    const onSubmit = async (e) => {
        e.preventDefault()
        setMsg('')
        try {
            await login({ email, password })
            setMsg('Login ok')
        } catch {
            setMsg('Credenciales inválidas')
        }
    }

    return (
        <form onSubmit={onSubmit}>
            <h2>Iniciar sesión</h2>
            <input placeholder="Email" value={email} onChange={e=>setEmail(e.target.value)} />
            <input type="password" placeholder="Password" value={password} onChange={e=>setPassword(e.target.value)} />
            <button type="submit">Entrar</button>
            {msg && <p>{msg}</p>}
        </form>
    )
}
