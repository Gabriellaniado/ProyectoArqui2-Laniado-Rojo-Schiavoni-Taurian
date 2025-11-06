import React, { useState, useEffect } from 'react';
import { searchService } from '../services/searchService';
import Header from '../components/Header';
import ProductCard from '../components/ProductCard';
import './ProductsPage.css';

const ProductsPage = () => {
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 9;

  useEffect(() => {
    fetchProducts();
  }, [currentPage]);

  const fetchProducts = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await searchService.searchProducts({
        page: currentPage,
        count: itemsPerPage
      });

      // La respuesta viene como { item: {...} }
      if (response.results && Array.isArray(response.results)) {
        setProducts(response.results);
      } else if (Array.isArray(response)) {
        setProducts(response);
      } else {
        setProducts([]);
      }
    } catch (err) {
      setError('Error al cargar los productos');
      console.error('Error fetching products:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = async (e) => {
    e.preventDefault();
    if (!searchQuery.trim()) {
      fetchProducts();
      return;
    }

    try {
      setLoading(true);
      setError(null);
      const response = await searchService.searchProducts({
        query: searchQuery,
        page: 1,
        count: itemsPerPage
      });

      if (response.results && Array.isArray(response.results)) {
        setProducts(response.results);
      } else if (Array.isArray(response)) {
        setProducts(response);
      } else {
        setProducts([]);
      }
      setCurrentPage(1);
    } catch (err) {
      setError('Error al buscar productos');
      console.error('Error searching products:', err);
    } finally {
      setLoading(false);
    }
  };

  const handlePreviousPage = () => {
    if (currentPage > 1) {
      setCurrentPage(currentPage - 1);
    }
  };

  const handleNextPage = () => {
    setCurrentPage(currentPage + 1);
  };

  return (
    <div className="products-page">
      <Header />

      <div className="container">
        <div className="search-section">
          <form onSubmit={handleSearch} className="search-form">
            <input
              type="text"
              className="search-input"
              placeholder="Buscar mates..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
            <button type="submit" className="search-button">
              Buscar
            </button>
          </form>
        </div>

        {loading && (
          <div className="loading">
            <p>Cargando productos...</p>
          </div>
        )}

        {error && (
          <div className="error-message">
            <p>{error}</p>
          </div>
        )}

        {!loading && !error && products.length === 0 && (
          <div className="no-products">
            <p>No se encontraron productos</p>
          </div>
        )}

        {!loading && !error && products.length > 0 && (
          <>
            <div className="products-grid">
              {products.map((product) => (
                <ProductCard key={product.id} product={product} />
              ))}
            </div>

            <div className="pagination">
              <button
                className="pagination-btn"
                onClick={handlePreviousPage}
                disabled={currentPage === 1}
              >
                ← Anterior
              </button>
              <span className="page-number">Página {currentPage}</span>
              <button
                className="pagination-btn"
                onClick={handleNextPage}
                disabled={products.length < itemsPerPage}
              >
                Siguiente →
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  );
};

export default ProductsPage;
