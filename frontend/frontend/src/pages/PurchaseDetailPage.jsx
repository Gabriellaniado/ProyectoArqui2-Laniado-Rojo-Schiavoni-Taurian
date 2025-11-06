import React from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import Header from '../components/Header';
import './PurchaseDetailPage.css';

const PurchaseDetailPage = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const purchase = location.state?.purchase;

  if (!purchase) {
    return (
      <div className="purchase-detail-page">
        <Header />
        <div className="container">
          <div className="error-message">
            No se encontró información de la compra
          </div>
          <button className="btn-back" onClick={() => navigate('/mis-compras')}>
            Volver a mis compras
          </button>
        </div>
      </div>
    );
  }

  const formatDate = (dateString) => {
    if (!dateString) return 'Fecha no disponible';
    const date = new Date(dateString);
    return date.toLocaleDateString('es-AR', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  return (
    <div className="purchase-detail-page">
      <Header />
      
      <div className="container">
        <button className="btn-back" onClick={() => navigate('/mis-compras')}>
          ← Volver a mis compras
        </button>

        <div className="purchase-detail-card">
          <h1 className="detail-title">Detalle de Compra</h1>

          <div className="detail-content">
            {purchase.productImage && (
              <div className="detail-image">
                <img 
                  src={purchase.productImage} 
                  alt={purchase.productName}
                  onError={(e) => {
                     e.currentTarget.style.display = 'none';
                  }}
                />
              </div>
            )}

            <div className="detail-info">
              <div className="info-section">
                <h2 className="section-title">Información del Producto</h2>
                
                <div className="info-row">
                  <span className="info-label">Producto:</span>
                  <span className="info-value">{purchase.productName}</span>
                </div>


              </div>

              <div className="info-section">
                <h2 className="section-title">Detalles de la Compra</h2>
                
                <div className="info-row">
                  <span className="info-label">Cantidad:</span>
                  <span className="info-value">{purchase.quantity} {purchase.quantity === 1 ? 'unidad' : 'unidades'}</span>
                </div>

                <div className="info-row">
                  <span className="info-label">Precio Total:</span>
                  <span className="info-value highlight">${purchase.total_price?.toFixed(2)}</span>
                </div>

                <div className="info-row">
                  <span className="info-label">Fecha de Compra:</span>
                  <span className="info-value">{formatDate(purchase.created_at || purchase.date)}</span>
                </div>
              </div>

              <div className="price-summary">
                <div className="summary-row">
                  <span>Precio por unidad:</span>
                  <span>${(purchase.total_price / purchase.quantity).toFixed(2)}</span>
                </div>
                <div className="summary-row">
                  <span>Cantidad:</span>
                  <span>x {purchase.quantity}</span>
                </div>
                <div className="summary-row total">
                  <span>Total Pagado:</span>
                  <span>${purchase.total_price?.toFixed(2)}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default PurchaseDetailPage;
