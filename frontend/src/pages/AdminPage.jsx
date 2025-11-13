import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { searchService } from '../services/searchService';
import { isAdmin } from '../utils/auth';
import Header from '../components/Header';
import AdminProductCard from '../components/AdminProductCard';
import './AdminPage.css';

const AdminPage = () => {
  const navigate = useNavigate();
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 8;

  // Filtros
  const [searchQuery, setSearchQuery] = useState('');
  const [minPrice, setMinPrice] = useState('');
  const [maxPrice, setMaxPrice] = useState('');
  const [appliedFilters, setAppliedFilters] = useState({
    name: '',
    minPrice: null,
    maxPrice: null,
  });

  // Verificar que sea admin
  useEffect(() => {
    if (!isAdmin()) {
      navigate('/');
    }
  }, [navigate]);

  // Cargar productos
  useEffect(() => {
    fetchProducts();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentPage, appliedFilters]);

  const fetchProducts = async () => {
    try {
      setLoading(true);
      setError(null);

      const filters = {
        page: currentPage,
        count: itemsPerPage,
      };

      if (appliedFilters.name) {
        filters.name = appliedFilters.name;
      }
      if (appliedFilters.minPrice && !isNaN(appliedFilters.minPrice)) {
        filters.minPrice = parseFloat(appliedFilters.minPrice);
      }
      if (appliedFilters.maxPrice && !isNaN(appliedFilters.maxPrice)) {
        filters.maxPrice = parseFloat(appliedFilters.maxPrice);
      }

      const response = await searchService.searchProducts(filters);

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

  const handleSearch = (e) => {
    e.preventDefault();
    setCurrentPage(1);
    setAppliedFilters({
      ...appliedFilters,
      name: searchQuery,
    });
  };

  const handleApplyPriceFilters = () => {
    const min = minPrice ? parseFloat(minPrice) : null;
    const max = maxPrice ? parseFloat(maxPrice) : null;

    if (min !== null && max !== null && max < min) {
      alert('El precio máximo no puede ser menor al precio mínimo');
      return;
    }

    setCurrentPage(1);
    setAppliedFilters({
      ...appliedFilters,
      minPrice: minPrice || null,
      maxPrice: maxPrice || null,
    });
  };

  const handleClearFilters = () => {
    setSearchQuery('');
    setMinPrice('');
    setMaxPrice('');
    setCurrentPage(1);
    setAppliedFilters({
      name: '',
      minPrice: null,
      maxPrice: null,
    });
  };

  const handleProductDelete = (productId) => {
    setProducts(products.filter(p => p.id !== productId));
  };

  const handlePreviousPage = () => {
    if (currentPage > 1) {
      setCurrentPage(currentPage - 1);
    }
  };

  const handleNextPage = () => {
    setCurrentPage(currentPage + 1);
  };

  const hasActiveFilters =
    appliedFilters.name || appliedFilters.minPrice || appliedFilters.maxPrice;

  return (
    <div className="admin-page">
      <Header />

      <div className="container">
        <div className="admin-header">
          <h1 className="admin-title">Administración de Productos</h1>
          <button 
            className="btn-create-product"
            onClick={() => navigate('/admin/nuevo-producto')}
          >
            + Crear Nuevo Producto
          </button>
        </div>

        {/* Barra de búsqueda */}
        <div className="search-section">
          <form onSubmit={handleSearch} className="search-form">
            <input
              type="text"
              className="search-input"
              placeholder="Buscar productos..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
            <button type="submit" className="search-button">
              Buscar
            </button>
          </form>
        </div>

        {/* Filtros de precio */}
        <div className="filters-section">
          <div className="price-filters">
            <div className="price-filter-group">
              <label htmlFor="minPrice" className="filter-label">
                Precio Mínimo
              </label>
              <input
                type="number"
                id="minPrice"
                className="price-input"
                placeholder="$ Min"
                min="0"
                step="100"
                value={minPrice}
                onChange={(e) => setMinPrice(e.target.value)}
              />
            </div>

            <div className="price-filter-separator">-</div>

            <div className="price-filter-group">
              <label htmlFor="maxPrice" className="filter-label">
                Precio Máximo
              </label>
              <input
                type="number"
                id="maxPrice"
                className="price-input"
                placeholder="$ Max"
                min="0"
                step="100"
                value={maxPrice}
                onChange={(e) => setMaxPrice(e.target.value)}
              />
            </div>

            <button
              type="button"
              className="apply-filters-btn"
              onClick={handleApplyPriceFilters}
            >
              Aplicar Filtros
            </button>

            {(minPrice || maxPrice || searchQuery || hasActiveFilters) && (
              <button
                type="button"
                className="clear-filters-btn"
                onClick={handleClearFilters}
              >
                Limpiar
              </button>
            )}
          </div>
        </div>

        {/* Contenido */}
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
            {hasActiveFilters && (
              <button className="clear-filters-btn" onClick={handleClearFilters}>
                Limpiar Filtros
              </button>
            )}
          </div>
        )}

        {!loading && !error && products.length > 0 && (
          <>
            <div className="products-grid">
              {products.map((product) => (
                <AdminProductCard 
                  key={product.id} 
                  product={product}
                  onDelete={handleProductDelete}
                />
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

export default AdminPage;
