import { Link } from 'react-router-dom'

export default function ProductCard({ product }) {
  const id = product.id || product.ID
  const name = product.name || product.Name
  const price = product.price ?? product.Price
  const img = product.image_url || product.ImageURL
  return (
    <article className="card product">
      <div className="thumb">
        {img ? (
          <img src={img} alt={name} />
        ) : (
          <div className="ph" aria-hidden />
        )}
      </div>
      <div className="info">
        <h3 className="title">{name}</h3>
        <p className="price">${price}</p>
        <Link className="btn btn-outline" to={`/product/${id}`}>Ver detalle</Link>
      </div>
    </article>
  )
}
