import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { salesService } from '../services/salesService';
import { productService } from '../services/productService';
import { getCustomerId, getCustomerIDFromToken } from '../utils/auth';
import { isAuthenticated } from '../utils/auth';
import Header from '../components/Header';
import './PurchasesPage.css';

const PurchasesPage = () => {
  const navigate = useNavigate();
  const [purchases, setPurchases] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    // Verificar autenticación
    if (!isAuthenticated()) {
      navigate('/login');
      return;
    }

    fetchPurchases();
  }, []);

  const fetchPurchases = async () => {
    try {
      setLoading(true);
      setError(null);
      const customerId = getCustomerIDFromToken();

      if (!customerId) {
        setError('No se pudo obtener el ID del usuario');
        return;
      }

      const response = await salesService.getSalesByCustomerId(customerId);

      // La respuesta viene como { customer_id, sales, count, total_spent }
      let salesData = [];
      if (response.sales && Array.isArray(response.sales)) {
        salesData = response.sales;
      } else if (Array.isArray(response)) {
        salesData = response;
      }

      // Obtener información de los productos para cada compra
      const purchasesWithProducts = await Promise.all(
        salesData.map(async (sale) => {
          try {
            const productResponse = await productService.getProductById(sale.item_id);
            const product = productResponse.item || productResponse;
            return {
              ...sale,
              productName: product.name,
              productImage: product.image_url
            };
          } catch (err) {
            console.error('Error fetching product:', err);
            return {
              ...sale,
              productName: 'Producto no disponible',
              productImage: null
            };
          }
        })
      );

      setPurchases(purchasesWithProducts);
    } catch (err) {
      console.error('Error fetching purchases:', err);
      setError('Error al cargar las compras');
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString) => {
    if (!dateString) return 'Fecha no disponible';
    const date = new Date(dateString);
    return date.toLocaleDateString('es-AR', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  const handleViewDetail = (purchase) => {
    navigate(`/compra/${purchase.id}`, { state: { purchase } });
  };

  if (loading) {
    return (
      <div className="purchases-page">
        <Header />
        <div className="container">
          <div className="loading">Cargando compras...</div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="purchases-page">
        <Header />
        <div className="container">
          <div className="error-message">{error}</div>
        </div>
      </div>
    );
  }

  return (
    <div className="purchases-page">
      <Header />

      <div className="container">
        <h1 className="page-title">Mis Compras</h1>

        {purchases.length === 0 ? (
          <div className="no-purchases">
            <p>Aún no tienes compras realizadas</p>
            <button className="btn-go-shopping" onClick={() => navigate('/')}>
              Ir a comprar
            </button>
          </div>
        ) : (
          <div className="purchases-list">
            {purchases.map((purchase, index) => (
              <div key={purchase.id || index} className="purchase-card">
                <div className="purchase-number">
                  Compra #{index + 1}
                </div>

                <div className="purchase-content">
                  {purchase.productImage && (
                    <div className="purchase-image">
                      <img
                        src={purchase.productImage}
                        alt={purchase.productName}
                        onError={(e) => {
                          e.currentTarget.style.display = 'none';
                        }}
                      />
                    </div>
                  )}

                  <div className="purchase-info">
                    <h3 className="purchase-product-name">{purchase.productName}</h3>
                    <p className="purchase-quantity">
                      <span className="label">Cantidad:</span> {purchase.quantity}
                    </p>
                    <p className="purchase-price">
                      <span className="label">Total:</span> ${purchase.total_price?.toFixed(2)}
                    </p>
                    <p className="purchase-date">
                      <span className="label">Fecha:</span> {formatDate(purchase.sale_date)}
                    </p>
                  </div>
                </div>

                <button
                  className="btn-detail"
                  onClick={() => handleViewDetail(purchase)}
                >
                  Ver detalle
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default PurchasesPage;
