export default function SearchBox({ value, onChange, onSearch }) {
    return (
        <div style={{ display:'flex', gap:8 }}>
            <input
                placeholder="Buscar…"
                value={value}
                onChange={(e)=>onChange(e.target.value)}
            />
            <button onClick={onSearch}>Buscar</button>
        </div>
    )
}
