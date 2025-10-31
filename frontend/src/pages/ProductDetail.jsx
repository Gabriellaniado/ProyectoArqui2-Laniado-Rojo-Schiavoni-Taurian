import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { getProductById } from '../services/products'

export default function ProductDetail() {
    const { id } = useParams()
    const [data, setData] = useState(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)

    useEffect(() => {
        (async () => {
            try {
                const res = await getProductById(id)
                setData(res)
            } catch (e) {
                setError('No pudimos cargar el producto')
            } finally {
                setLoading(false)
            }
        })()
    }, [id])

    if (loading) return <p>Cargando…</p>
    if (error) return <p>{error}</p>
    if (!data) return null

    return (
        <article>
            <h2>{data.name}</h2>
            <p>Precio: ${data.price}</p>
            <button>Agregar al carrito</button>
        </article>
    )
}
