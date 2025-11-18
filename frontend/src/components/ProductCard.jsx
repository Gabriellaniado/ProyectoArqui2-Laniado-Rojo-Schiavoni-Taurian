
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCart } from '../context/CartContext';
import { isAuthenticated } from '../utils/auth';
import './ProductCard.css';

const ProductCard = ({ product, isAdmin = false, onEdit, onDelete }) => {
    const navigate = useNavigate();
    const { addItem } = useCart();
    const [adding, setAdding] = useState(false);

    const handleViewDetails = () => {
        navigate(`/producto/${product.id}`);
    };

    const handleAddToCart = async (e) => {
        e.stopPropagation(); // Evitar que se propague el click

        // Verificar autenticación
        if (!isAuthenticated()) {
            alert('Debes iniciar sesión para agregar productos al carrito');
            navigate('/login');
            return;
        }

        try {
            setAdding(true);
            const success = await addItem(product.id, 1);
            if (success) {
                // Pequeña notificación de éxito
                alert('✅ Producto agregado al carrito');
            }
        } catch (err) {
            console.error('Error adding to cart:', err);
        } finally {
            setAdding(false);
        }
    };

    const handleEdit = (e) => {
        e.stopPropagation();
        if (onEdit) {
            onEdit(product);
        }
    };

    const handleDelete = (e) => {
        e.stopPropagation();
        if (onDelete) {
            const confirmDelete = window.confirm(
                `¿Estás seguro de que deseas eliminar el producto "${product.name}"?`
            );
            if (confirmDelete) {
                onDelete(product.id);
            }
        }
    };

    return (
        <div className="product-card">
            <div className="product-image" onClick={handleViewDetails}>
                <img
                    src={product.image_url}
                    alt={product.name}
                    onError={(e) => {
                        e.currentTarget.style.display = 'none';
                    }}
                />
            </div>
            <div className="product-info">
                <h3 className="product-name">{product.name}</h3>
                <p className="product-price">${product.price.toFixed(2)}</p>

                {isAdmin ? (
                    // Vista de administrador
                    <div className="product-admin-actions">
                        <button
                            className="btn-edit-product"
                            onClick={handleEdit}
                            title="Editar producto"
                        >
                            ✏️ Editar
                        </button>
                        <button
                            className="btn-delete-product"
                            onClick={handleDelete}
                            title="Eliminar producto"
                        >
                            🗑️ Eliminar
                        </button>
                    </div>
                ) : (
                    // Vista de cliente
                    <div className="product-card-actions">
                        <button
                            className="btn-add-to-cart-card"
                            onClick={handleAddToCart}
                            disabled={adding || product.stock === 0}
                            title="Agregar al carrito"
                        >
                            {adding ? '...' : '🛒'}
                        </button>
                        <button className="btn-view-details" onClick={handleViewDetails}>
                            Ver Detalles
                        </button>
                    </div>
                )}
            </div>
        </div>
    );
};

export default ProductCard;