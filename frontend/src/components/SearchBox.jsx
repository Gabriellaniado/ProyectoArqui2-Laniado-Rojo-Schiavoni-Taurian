export default function SearchBox({ value, onChange, onSearch }) {
  const handleKey = (e) => {
    if (e.key === 'Enter') onSearch()
  }
  return (
    <div className="searchbox">
      <input
        className="input"
        placeholder="Buscar…"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onKeyDown={handleKey}
      />
      <button className="btn btn-primary" onClick={onSearch}>Buscar</button>
    </div>
  )
}

