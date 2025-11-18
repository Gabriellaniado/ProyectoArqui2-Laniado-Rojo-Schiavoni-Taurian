import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useCart } from '../context/CartContext';
import './CartDrawer.css';

const CartDrawer = () => {
    const navigate = useNavigate();
    const {
        cart,
        loading,
        isOpen,
        closeCart,
        updateItem,
        removeItem,
    } = useCart();

    const handleQuantityChange = async (itemID, newQuantity) => {
        if (newQuantity < 1) return;
        await updateItem(itemID, newQuantity);
    };

    const handleRemove = async (itemID) => {
        if (window.confirm('¿Eliminar este producto del carrito?')) {
            await removeItem(itemID);
        }
    };

    const handleGoToCart = () => {
        closeCart();
        navigate('/carrito');
    };

    if (!isOpen) return null;

    return (
        <>
            {/* Overlay oscuro */}
            <div className="cart-overlay" onClick={closeCart}></div>

            {/* Drawer del carrito */}
            <div className={`cart-drawer ${isOpen ? 'open' : ''}`}>
                <div className="cart-drawer-header">
                    <h2>🛒 Mi Carrito</h2>
                    <button className="close-btn" onClick={closeCart}>
                        ✕
                    </button>
                </div>

                <div className="cart-drawer-content">
                    {loading && <div className="cart-loading">Cargando...</div>}

                    {!loading && cart.items.length === 0 && (
                        <div className="cart-empty">
                            <p>Tu carrito está vacío</p>
                            <button className="btn-continue" onClick={closeCart}>
                                Continuar comprando
                            </button>
                        </div>
                    )}

                    {!loading && cart.items.length > 0 && (
                        <>
                            <div className="cart-items">
                                {cart.items.map((item) => (
                                    <div key={item.item_id} className="cart-item">
                                        <img
                                            src={item.image_url}
                                            alt={item.name}
                                            className="cart-item-image"
                                        />
                                        <div className="cart-item-info">
                                            <h4>{item.name}</h4>
                                            <p className="cart-item-price">${item.price.toFixed(2)}</p>

                                            <div className="cart-item-quantity">
                                                <button
                                                    className="qty-btn"
                                                    onClick={() => handleQuantityChange(item.item_id, item.quantity - 1)}
                                                    disabled={item.quantity <= 1}
                                                >
                                                    -
                                                </button>
                                                <span className="qty-value">{item.quantity}</span>
                                                <button
                                                    className="qty-btn"
                                                    onClick={() => handleQuantityChange(item.item_id, item.quantity + 1)}
                                                    disabled={item.quantity >= item.stock}
                                                >
                                                    +
                                                </button>
                                            </div>

                                            <p className="cart-item-subtotal">
                                                Subtotal: ${item.subtotal.toFixed(2)}
                                            </p>
                                        </div>
                                        <button
                                            className="remove-btn"
                                            onClick={() => handleRemove(item.item_id)}
                                            title="Eliminar producto"
                                        >
                                            🗑️
                                        </button>
                                    </div>
                                ))}
                            </div>

                            <div className="cart-drawer-footer">
                                <div className="cart-total">
                                    <span>Total:</span>
                                    <span className="total-amount">${cart.total.toFixed(2)}</span>
                                </div>
                                <button className="btn-checkout" onClick={handleGoToCart}>
                                    Ver Carrito Completo
                                </button>
                            </div>
                        </>
                    )}
                </div>
            </div>
        </>
    );
};

export default CartDrawer;