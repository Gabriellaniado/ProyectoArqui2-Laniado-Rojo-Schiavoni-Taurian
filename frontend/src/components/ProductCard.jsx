export default function ProductCard({ product }) {
    return (
        <article className="card">
            <h3>{product.name}</h3>
            <p>${product.price}</p>
            <a href={`/product/${product.id}`}>Ver detalle</a>
        </article>
    )
}
