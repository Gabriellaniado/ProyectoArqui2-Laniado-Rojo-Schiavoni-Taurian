import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, useLocation } from 'react-router-dom';
import { productService } from '../services/productService';
import { salesService } from '../services/salesService';
import { isAuthenticated, getCustomerId } from '../utils/auth';
import Header from '../components/Header';
import './ProductDetailPage.css';
import { useCart } from '../context/CartContext';

const ProductDetailPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const [product, setProduct] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [quantity, setQuantity] = useState(1);
  const [purchasing, setPurchasing] = useState(false);
  const { addItem } = useCart();
  const [addingToCart, setAddingToCart] = useState(false);

  const savedFilters = location.state?.filters;

  useEffect(() => {
    fetchProduct();
  }, [id]);

  const fetchProduct = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await productService.getProductById(id);

      // La respuesta viene como { item: {...} }
      if (response.item) {
        setProduct(response.item);
      } else {
        setProduct(response);
      }
    } catch (err) {
      setError('Error al cargar el producto');
      console.error('Error fetching product:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleQuantityChange = (e) => {
    const value = parseInt(e.target.value);
    if (!isNaN(value) && value > 0) {
      setQuantity(value);
    }
  };

  const handleAddToCart = async () => {
    // Verificar autenticación
    if (!isAuthenticated()) {
      alert('Debes iniciar sesión para agregar productos al carrito');
      navigate('/login');
      return;
    }

    try {
      setAddingToCart(true);
      const success = await addItem(product.id, quantity);
      if (success) {
        alert('✅ Producto agregado al carrito');
        setQuantity(1); // Resetear cantidad
      }
    } catch (err) {
      console.error('Error adding to cart:', err);
    } finally {
      setAddingToCart(false);
    }
  };

  const handlePurchase = async () => {
    // Verificar autenticación
    if (!isAuthenticated()) {
      alert('Debes iniciar sesión para realizar una compra');
      navigate('/login');
      return;
    }

    // Confirmar compra
    const confirmPurchase = window.confirm(
      `¿Estás seguro de comprar la cantidad ${quantity} solicitada de ${product.name}?`
    );

    if (!confirmPurchase) return;

    try {
      setPurchasing(true);
      const customerId = (getCustomerId());

      const saleData = {
        item_id: product.id,
        quantity: quantity,
        customer_id: customerId
      };

      const response = await salesService.createSale(saleData);

      // Si el status es 201, mostrar confirmación
      if (response) {
        alert('Compra confirmada');
        navigate('/mis-compras');
      }
    } catch (err) {
      console.error('Error creating sale:', err);
      alert('Error al procesar la compra. Por favor intenta nuevamente.');
    } finally {
      setPurchasing(false);
    }
  };

    const handleGoBack = () => {
  console.log('🔙 Volviendo atrás con filtros:', savedFilters); // 👈 Debug
  navigate('/', { state: savedFilters ? { filters: savedFilters } : undefined });
};

  if (loading) {
    return (
      <div className="product-detail-page">
        <Header />
        <div className="container">
          <div className="loading">Cargando producto...</div>
        </div>
      </div>
    );
  }

  if (error || !product) {
    return (
      <div className="product-detail-page">
        <Header />
        <div className="container">
          <div className="error-message">{error || 'Producto no encontrado'}</div>
        </div>
      </div>
    );
  }

  return (
    <div className="product-detail-page">
      <Header />

      <div className="container">
        <button className="back-button" onClick={handleGoBack}>
          ← Volver a productos
        </button>

        <div className="product-detail">
          <div className="product-image-large">
            <img
              src={product.image_url}
              alt={product.name}
              onError={(e) => {
                
              }}
            />
          </div>

          <div className="product-info-detailed">
            <h1 className="product-title">{product.name}</h1>

            <div className="product-category">
              <span className="category-label">Categoría:</span>
              <span className="category-value">{product.category}</span>
            </div>

            <div className="product-description">
              <h3>Descripción</h3>
              <p>{product.description}</p>
            </div>

            <div className="product-stock">
              <span className="stock-label">Stock disponible:</span>
              <span className="stock-value">{product.stock} unidades</span>
            </div>

            <div className="product-price-large">
              ${product.price.toFixed(2)}
            </div>

            <div className="purchase-section">
              <div className="quantity-selector">
                <label htmlFor="quantity">Cantidad:</label>
                <input
                    type="number"
                    id="quantity"
                    min="1"
                    max={product.stock}
                    value={quantity}
                    onChange={handleQuantityChange}
                    className="quantity-input"
                />
              </div>

              <div className="total-price">
                <span>Total:</span>
                <span className="total-amount">
                ${(product.price * quantity).toFixed(2)}
                </span>
              </div>

              <div className="action-buttons">
                <button
                    className="btn-add-to-cart"
                    onClick={handleAddToCart}
                    disabled={addingToCart || quantity > product.stock}
                >
                  {addingToCart ? 'Agregando...' : '🛒 Agregar al Carrito'}
                </button>

                <button
                    className="btn-purchase"
                    onClick={handlePurchase}
                    disabled={purchasing || quantity > product.stock}
                >
                  {purchasing ? 'Procesando...' : 'Comprar Ahora'}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProductDetailPage;
