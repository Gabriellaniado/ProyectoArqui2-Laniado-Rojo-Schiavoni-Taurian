import { useEffect, useState } from 'react'
import { listProducts } from '../services/products'
import SearchBox from '../components/SearchBox.jsx'
import ProductCard from '../components/ProductCard.jsx'

export default function Catalog() {
    const [q, setQ] = useState('')
    const [items, setItems] = useState([])
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState(null)

    const fetchData = async () => {
        try {
            setLoading(true); setError(null)
            const data = await listProducts({ q })
            setItems(data?.items ?? [])
        } catch (e) {
            setError('No pudimos cargar el catálogo')
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => { fetchData() }, []) // carga inicial

    return (
        <section>
            <h2>Catálogo</h2>
            <SearchBox value={q} onChange={setQ} onSearch={fetchData} />
            {loading && <p>Cargando…</p>}
            {error && <p>{error}</p>}
            <div className="grid">
                {items.map(p => <ProductCard key={p.id} product={p} />)}
            </div>
        </section>
    )
}
