import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCart } from '../context/CartContext';
import Header from '../components/Header';
import './CartPage.css';

const CartPage = () => {
    const navigate = useNavigate();
    const {
        cart,
        loading,
        updateItem,
        removeItem,
        clearCart,
        checkout,
    } = useCart();

    const [processingCheckout, setProcessingCheckout] = useState(false);

    const handleQuantityChange = async (itemID, newQuantity) => {
        if (newQuantity < 1) return;
        await updateItem(itemID, newQuantity);
    };

    const handleRemove = async (itemID) => {
        if (window.confirm('¿Eliminar este producto del carrito?')) {
            await removeItem(itemID);
        }
    };

    const handleClearCart = async () => {
        if (window.confirm('¿Estás seguro de vaciar todo el carrito?')) {
            await clearCart();
        }
    };

    const handleCheckout = async () => {
        if (cart.items.length === 0) {
            alert('El carrito está vacío');
            return;
        }

        // Verificar stock antes de procesar
        const insufficientStock = cart.items.find(item => item.quantity > item.stock);
        if (insufficientStock) {
            alert(`Stock insuficiente para ${insufficientStock.name}. Disponible: ${insufficientStock.stock}`);
            return;
        }

        const confirmPurchase = window.confirm(
            `¿Confirmar compra?\n\nTotal: $${cart.total.toFixed(2)}\nProductos: ${cart.item_count}`
        );

        if (!confirmPurchase) return;

        try {
            setProcessingCheckout(true);
            await checkout();
            alert('¡Compra realizada con éxito! ✅');
            navigate('/mis-compras');
        } catch (error) {
            alert(error.message || 'Error al procesar la compra');
        } finally {
            setProcessingCheckout(false);
        }
    };

    return (
        <div className="cart-page">
            <Header />

            <div className="container">
                <div className="cart-page-header">
                    <h1>🛒 Mi Carrito</h1>
                    {cart.items.length > 0 && (
                        <button className="btn-clear-cart" onClick={handleClearCart}>
                            Vaciar Carrito
                        </button>
                    )}
                </div>

                {loading && (
                    <div className="cart-page-loading">
                        <p>Cargando carrito...</p>
                    </div>
                )}

                {!loading && cart.items.length === 0 && (
                    <div className="cart-page-empty">
                        <div className="empty-icon">🛒</div>
                        <h2>Tu carrito está vacío</h2>
                        <p>¡Agrega productos para comenzar a comprar!</p>
                        <button className="btn-continue-shopping" onClick={() => navigate('/')}>
                            Ver Productos
                        </button>
                    </div>
                )}

                {!loading && cart.items.length > 0 && (
                    <div className="cart-page-content">
                        <div className="cart-items-section">
                            <h2>Productos ({cart.item_count})</h2>

                            {cart.items.map((item) => (
                                <div key={item.item_id} className="cart-page-item">
                                    <img
                                        src={item.image_url}
                                        alt={item.name}
                                        className="cart-page-item-image"
                                    />

                                    <div className="cart-page-item-details">
                                        <h3>{item.name}</h3>
                                        <p className="item-description">{item.description}</p>
                                        <p className="item-stock">Stock disponible: {item.stock} unidades</p>
                                    </div>

                                    <div className="cart-page-item-price">
                                        <p className="price-label">Precio unitario</p>
                                        <p className="price-value">${item.price.toFixed(2)}</p>
                                    </div>

                                    <div className="cart-page-item-quantity">
                                        <p className="quantity-label">Cantidad</p>
                                        <div className="quantity-controls">
                                            <button
                                                className="qty-btn"
                                                onClick={() => handleQuantityChange(item.item_id, item.quantity - 1)}
                                                disabled={item.quantity <= 1}
                                            >
                                                -
                                            </button>
                                            <input
                                                type="number"
                                                className="qty-input"
                                                value={item.quantity}
                                                onChange={(e) => {
                                                    const value = parseInt(e.target.value);
                                                    if (!isNaN(value) && value > 0 && value <= item.stock) {
                                                        handleQuantityChange(item.item_id, value);
                                                    }
                                                }}
                                                min="1"
                                                max={item.stock}
                                            />
                                            <button
                                                className="qty-btn"
                                                onClick={() => handleQuantityChange(item.item_id, item.quantity + 1)}
                                                disabled={item.quantity >= item.stock}
                                            >
                                                +
                                            </button>
                                        </div>
                                    </div>

                                    <div className="cart-page-item-subtotal">
                                        <p className="subtotal-label">Subtotal</p>
                                        <p className="subtotal-value">${item.subtotal.toFixed(2)}</p>
                                    </div>

                                    <button
                                        className="btn-remove-item"
                                        onClick={() => handleRemove(item.item_id)}
                                        title="Eliminar producto"
                                    >
                                        🗑️
                                    </button>
                                </div>
                            ))}
                        </div>

                        <div className="cart-summary-section">
                            <div className="cart-summary">
                                <h2>Resumen de Compra</h2>

                                <div className="summary-row">
                                    <span>Productos ({cart.item_count})</span>
                                    <span>${cart.total.toFixed(2)}</span>
                                </div>

                                <div className="summary-row summary-shipping">
                                    <span>Envío</span>
                                    <span className="free-shipping">Gratis</span>
                                </div>

                                <div className="summary-divider"></div>

                                <div className="summary-row summary-total">
                                    <span>Total</span>
                                    <span className="total-amount">${cart.total.toFixed(2)}</span>
                                </div>

                                <button
                                    className="btn-checkout-main"
                                    onClick={handleCheckout}
                                    disabled={processingCheckout}
                                >
                                    {processingCheckout ? 'Procesando...' : 'Finalizar Compra'}
                                </button>

                                <button
                                    className="btn-continue-shopping-secondary"
                                    onClick={() => navigate('/')}
                                >
                                    Continuar Comprando
                                </button>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};

export default CartPage;