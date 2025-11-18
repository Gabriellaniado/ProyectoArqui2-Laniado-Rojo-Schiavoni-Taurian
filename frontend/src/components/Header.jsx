import React from 'react';
import { useNavigate } from 'react-router-dom';
import { isAuthenticated, isAdmin, removeToken } from '../utils/auth';
import { useCart } from '../context/CartContext';
import './Header.css';

const Header = () => {
  const navigate = useNavigate();
  const authenticated = isAuthenticated();
  const userIsAdmin = isAdmin();
  const { cart, toggleCart} = useCart();

  const handleLogout = () => {
    removeToken();
    navigate('/');
  };

  return (
      <header className="header">
        <div className="header-container">
          <div className="logo" onClick={() => navigate('/')}>
            <img
                src="/logo-gustoamate.jpg"
                alt="GustoaMate"
                className="logo-image"
            />
            <span className="logo-text">Gusto a Mate</span>
          </div>
          <nav className="nav-buttons">
            {authenticated && (
                <button
                    className="btn-cart"
                    onClick={toggleCart}
                    title="Ver carrito"
                >
                  🛒
                  {cart.item_count > 0 && (
                      <span className="cart-badge">{cart.item_count}</span>
                  )}
                </button>
            )}

            {authenticated ? (
                <>
                  {userIsAdmin && (
                      <button
                          className="btn-secondary"
                          onClick={() => navigate('/admin')}
                      >
                        Panel Admin
                      </button>
                  )}
                  <button
                      className="btn-secondary"
                      onClick={() => navigate('/mis-compras')}
                  >
                    Mis Compras
                  </button>
                  <button
                      className="btn-primary"
                      onClick={handleLogout}
                  >
                    Cerrar Sesión
                  </button>
                </>
            ) : (
                <>
                  <button
                      className="btn-secondary"
                      onClick={() => navigate('/login')}
                  >
                    Iniciar Sesión
                  </button>
                  <button
                      className="btn-primary"
                      onClick={() => navigate('/registro')}
                  >
                    Registrarse
                  </button>
                </>
            )}
          </nav>
        </div>
      </header>
  );
};

export default Header;
