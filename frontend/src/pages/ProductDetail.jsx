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

  if (loading) return <p className="muted">Cargando…</p>
  if (error) return <p className="error">{error}</p>
  if (!data) return null

  return (
    <article className="product-detail">
      <div className="media">
        {data.image_url && <img src={data.image_url} alt={data.name} />}
      </div>
      <div className="summary">
        <h2>{data.name}</h2>
        <p className="price">${data.price}</p>
        <p className="desc">{data.description}</p>
        <button className="btn btn-primary">Agregar al carrito</button>
      </div>
    </article>
  )
}

