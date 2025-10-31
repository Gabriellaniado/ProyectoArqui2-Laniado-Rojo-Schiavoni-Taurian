import { Link, NavLink } from 'react-router-dom'

export default function Header() {
  return (
    <header className="header">
      <div className="container header-row">
        <Link to="/" className="brand">Mates Store</Link>
        <nav className="nav">
          <NavLink to="/catalog" className={({isActive}) => `nav-link${isActive ? ' active' : ''}`}>Catálogo</NavLink>
          <NavLink to="/cart" className={({isActive}) => `nav-link${isActive ? ' active' : ''}`}>Carrito</NavLink>
        </nav>
        <div className="spacer" />
        <nav className="nav right">
          <NavLink to="/login" className={({isActive}) => `btn btn-ghost${isActive ? ' active' : ''}`}>Login</NavLink>
        </nav>
      </div>
    </header>
  )
}
