import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { productService } from '../services/productService';
import { isAdmin } from '../utils/auth';
import Header from '../components/Header';
import './ProductFormPage.css';

const NewProductPage = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: '',
    category: '',
    description: '',
    price: '',
    stock: '',
    image_url: ''
  });
  const [errors, setErrors] = useState({});
  const [loading, setLoading] = useState(false);

  // Verificar que sea admin
  useEffect(() => {
    if (!isAdmin()) {
      navigate('/');
    }
  }, [navigate]);

  const validateURL = (url) => {
    try {
      new URL(url);
      return true;
    } catch {
      return false;
    }
  };

  const validateForm = () => {
    const newErrors = {};

    // Validar campos obligatorios
    if (!formData.name.trim()) {
      newErrors.name = 'El nombre es obligatorio';
    }

    if (!formData.category.trim()) {
      newErrors.category = 'La categoría es obligatoria';
    }

    if (!formData.description.trim()) {
      newErrors.description = 'La descripción es obligatoria';
    }

    // Validar precio
    if (!formData.price) {
      newErrors.price = 'El precio es obligatorio';
    } else if (isNaN(formData.price) || parseFloat(formData.price) <= 0) {
      newErrors.price = 'El precio debe ser un número mayor a 0';
    }

    // Validar stock
    if (formData.stock === '') {
      newErrors.stock = 'El stock es obligatorio';
    } else if (isNaN(formData.stock) || parseInt(formData.stock) < 0) {
      newErrors.stock = 'El stock debe ser un número mayor o igual a 0';
    }

    // Validar URL de imagen
    if (!formData.image_url.trim()) {
      newErrors.image_url = 'La URL de la imagen es obligatoria';
    } else if (!validateURL(formData.image_url)) {
      newErrors.image_url = 'La URL de la imagen no es válida';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });
    // Limpiar error del campo al escribir
    if (errors[name]) {
      setErrors({
        ...errors,
        [name]: ''
      });
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    const confirmCreate = window.confirm(
      `¿Estás seguro de crear el producto "${formData.name}"?`
    );

    if (!confirmCreate) return;

    try {
      setLoading(true);
      const productData = {
        name: formData.name,
        category: formData.category,
        description: formData.description,
        price: parseFloat(formData.price),
        stock: parseInt(formData.stock),
        image_url: formData.image_url
      };

      await productService.createProduct(productData);
      alert('Producto creado correctamente');
      navigate('/admin');
    } catch (error) {
      console.error('Error al crear producto:', error);
      alert('Error al crear el producto. Por favor intenta nuevamente.');
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    const confirmCancel = window.confirm(
      '¿Estás seguro de cancelar? Los cambios no se guardarán.'
    );

    if (confirmCancel) {
      navigate('/admin');
    }
  };

  return (
    <div className="product-form-page">
      <Header />

      <div className="container">
        <div className="form-header">
          <button className="btn-back" onClick={() => navigate('/admin')}>
            ← Volver
          </button>
          <h1 className="form-title">Crear Nuevo Producto</h1>
        </div>

        <form className="product-form" onSubmit={handleSubmit}>
          <div className="form-grid">
            <div className="form-group">
              <label htmlFor="name">Nombre del Producto *</label>
              <input
                type="text"
                id="name"
                name="name"
                value={formData.name}
                onChange={handleChange}
                className={errors.name ? 'error' : ''}
              />
              {errors.name && <span className="error-message">{errors.name}</span>}
            </div>

            <div className="form-group">
              <label htmlFor="category">Categoría *</label>
              <input
                type="text"
                id="category"
                name="category"
                value={formData.category}
                onChange={handleChange}
                className={errors.category ? 'error' : ''}
              />
              {errors.category && <span className="error-message">{errors.category}</span>}
            </div>

            <div className="form-group full-width">
              <label htmlFor="description">Descripción *</label>
              <textarea
                id="description"
                name="description"
                value={formData.description}
                onChange={handleChange}
                rows="4"
                className={errors.description ? 'error' : ''}
              />
              {errors.description && <span className="error-message">{errors.description}</span>}
            </div>

            <div className="form-group">
              <label htmlFor="price">Precio *</label>
              <input
                type="number"
                id="price"
                name="price"
                value={formData.price}
                onChange={handleChange}
                step="0.01"
                min="0"
                className={errors.price ? 'error' : ''}
              />
              {errors.price && <span className="error-message">{errors.price}</span>}
            </div>

            <div className="form-group">
              <label htmlFor="stock">Stock *</label>
              <input
                type="number"
                id="stock"
                name="stock"
                value={formData.stock}
                onChange={handleChange}
                min="0"
                className={errors.stock ? 'error' : ''}
              />
              {errors.stock && <span className="error-message">{errors.stock}</span>}
            </div>

            <div className="form-group full-width">
              <label htmlFor="image_url">URL de la Imagen *</label>
              <input
                type="text"
                id="image_url"
                name="image_url"
                value={formData.image_url}
                onChange={handleChange}
                placeholder="https://ejemplo.com/imagen.jpg"
                className={errors.image_url ? 'error' : ''}
              />
              {errors.image_url && <span className="error-message">{errors.image_url}</span>}
            </div>
          </div>

          <div className="form-actions">
            <button 
              type="button" 
              className="btn-cancel"
              onClick={handleCancel}
              disabled={loading}
            >
              Cancelar
            </button>
            <button 
              type="submit" 
              className="btn-submit"
              disabled={loading}
            >
              {loading ? 'Creando...' : 'Crear Producto'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default NewProductPage;
