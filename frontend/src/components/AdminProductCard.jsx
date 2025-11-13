import React from 'react';
import { useNavigate } from 'react-router-dom';
import { productService } from '../services/productService';
import './AdminProductCard.css';

const AdminProductCard = ({ product, onDelete }) => {
  const navigate = useNavigate();

  const handleEdit = () => {
    navigate(`/admin/editar-producto/${product.id}`);
  };

  const handleDelete = async () => {
    const confirmDelete = window.confirm(
      `¿Estás seguro de eliminar el producto "${product.name}"?`
    );

    if (!confirmDelete) return;

    try {
      await productService.deleteProduct(product.id);
      alert('Producto eliminado correctamente');
      onDelete(product.id);
    } catch (error) {
      console.error('Error al eliminar producto:', error);
      alert('Error al eliminar el producto. Por favor intenta nuevamente.');
    }
  };

  return (
    <div className="admin-product-card">
      <div className="product-image">
        <img 
          src={product.image_url || 'https://via.placeholder.com/300x400?text=Mate'} 
          alt={product.name}
          onError={(e) => {
            e.target.src = 'https://via.placeholder.com/300x400?text=Mate';
          }}
        />
      </div>
      <div className="product-info">
        <h3 className="product-name">{product.name}</h3>
        <p className="product-category">{product.category}</p>
        <p className="product-price">${product.price.toFixed(2)}</p>
        <p className="product-stock">Stock: {product.stock}</p>
        
        <div className="admin-buttons">
          <button className="btn-edit" onClick={handleEdit}>
            Modificar
          </button>
          <button className="btn-delete" onClick={handleDelete}>
            Eliminar
          </button>
        </div>
      </div>
    </div>
  );
};

export default AdminProductCard;
