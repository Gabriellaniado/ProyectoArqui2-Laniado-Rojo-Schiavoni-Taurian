import React, { useState, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { searchService } from '../services/searchService';
import Header from '../components/Header';
import ProductCard from '../components/ProductCard';
import './ProductsPage.css';

const ProductsPage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  // --- Estado ---
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 8;

  // Restaurar filtros desde location.state si existen
  const savedFilters = location.state?.filters;
  const savedSearchQuery = location.state?.searchQuery;
  // Estado para los *inputs* (lo que el usuario escribe)
  const [searchQuery, setSearchQuery] = useState('');
  const [minPrice, setMinPrice] = useState('');
  const [maxPrice, setMaxPrice] = useState('');

  // Estado para los *filtros aplicados* (lo que dispara la búsqueda)
  const [appliedFilters, setAppliedFilters] = useState({
    name: savedFilters?.name || '',
    minPrice: savedFilters?.minPrice || null,
    maxPrice: savedFilters?.maxPrice || null,
  });

  // --- Efecto Principal para Cargar Datos ---
  // Este useEffect es AHORA la *única* fuente de verdad para llamar a la API.
  // Se ejecuta si 'currentPage' o 'appliedFilters' cambian.
  useEffect(() => {
    const fetchProducts = async () => {
      try {
        setLoading(true);
        setError(null);

        // Construir los filtros solo con valores válidos
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

    fetchProducts();
    
    // Ya no necesitamos deshabilitar la regla de lint,
    // todas las dependencias están correctamente declaradas.
  }, [currentPage, appliedFilters]);

  // --- Manejadores de Eventos (Solo actualizan estado) ---

  /*const handleProductClick = (productId) => {
    console.log('📦 Navegando con filtros:', appliedFilters);
    navigate(`/products/${productId}`, {
      state: {
        filters: {
          name: appliedFilters.name,
          minPrice: appliedFilters.minPrice,
          maxPrice: appliedFilters.maxPrice,
        }
      }
    });
  };*/

  // Maneja el envío del formulario de búsqueda
  const handleSearch = (e) => {
    e.preventDefault();
    setCurrentPage(1); // Resetea la página
    setAppliedFilters({ // Aplica los filtros
      ...appliedFilters,
      name: searchQuery,
    });
  };

  // Maneja el clic en "Aplicar Filtros" de precio
  const handleApplyPriceFilters = () => {
    // Validación: solo cuando AMBOS campos tienen valor
    const min = minPrice ? parseFloat(minPrice) : null;
    const max = maxPrice ? parseFloat(maxPrice) : null;

    // Solo validar si ambos tienen valor
    if (min !== null && max !== null && max < min) {
      alert('El precio máximo no puede ser menor al precio mínimo');
      return;
    }

    setCurrentPage(1); // Resetea la página
    setAppliedFilters({ // Aplica los filtros
      ...appliedFilters,
      minPrice: minPrice || null, // Guarda null si está vacío
      maxPrice: maxPrice || null,
    });
  };

  // Limpia todos los inputs y filtros aplicados
  const handleClearFilters = () => {
    // 1. Limpiar los inputs
    setSearchQuery('');
    setMinPrice('');
    setMaxPrice('');
    // 2. Resetear la página
    setCurrentPage(1);
    // 3. Resetear los filtros aplicados (esto dispara el useEffect)
    setAppliedFilters({
      name: '',
      minPrice: null,
      maxPrice: null,
    });
  };

  // --- Paginación ---
  const handlePreviousPage = () => {
    if (currentPage > 1) {
      setCurrentPage(currentPage - 1);
    }
  };

  const handleNextPage = () => {
    setCurrentPage(currentPage + 1);
  };

  // Variable para saber si hay algún filtro activo
  const hasActiveFilters =
    appliedFilters.name || appliedFilters.minPrice || appliedFilters.maxPrice;

  return (
    <div className="products-page">
      <Header />

      <div className="container">
        {/* Barra de búsqueda - SEPARADA */}
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

        {/* Filtros de precio - SEPARADOS */}
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

        {/* --- Contenido Principal (Grid de Productos) --- */}

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
            {/* Mostrar botón de limpiar solo si había filtros aplicados */}
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