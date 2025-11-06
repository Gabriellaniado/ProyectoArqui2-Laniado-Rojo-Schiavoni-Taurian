import React from 'react';
import { useNavigate } from 'react-router-dom';
import './ProductCard.css';

const ProductCard = ({ product }) => {
  const navigate = useNavigate();

  const handleBuyClick = () => {
    navigate(`/producto/${product.id}`);
  };

  return (
    <div className="product-card">
      <div className="product-image">
        <img
          src={product.image_url}
          alt={product.name}
          onError={(e) => {
            e.target.src = 'https://via.placeholder.com/300x400?text=Mate';
          }}
        />
      </div>
      <div className="product-info">
        <h3 className="product-name">{product.name}</h3>
        <p className="product-price">${product.price.toFixed(2)}</p>
        <button className="btn-buy" onClick={handleBuyClick}>
          Comprar
        </button>
      </div>
    </div>
  );
};

export default ProductCard;
