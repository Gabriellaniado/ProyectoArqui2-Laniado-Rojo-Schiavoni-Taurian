import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { userService } from '../services/userService';
import './RegisterPage.css';

const RegisterPage = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    first_name: '',
    last_name: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);

    // Validar que todos los campos estén completos
    if (!formData.email || !formData.password || !formData.first_name || !formData.last_name) {
      setError('Todos los campos son obligatorios');
      return;
    }

    // Confirmar registro
    const confirmRegister = window.confirm('¿Estás seguro de registrarte con estos datos?');
    if (!confirmRegister) return;

    try {
      setLoading(true);
      const response = await userService.register(formData);
      
      // Si el registro fue exitoso (status 201)
      alert('Registro realizado correctamente. Ahora inicia sesión.');
      navigate('/login');
    } catch (err) {
      console.error('Error during registration:', err);
      setError('Error al registrarse. Por favor intenta nuevamente.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="register-page">
      <div className="register-container">
        <div className="register-header">
          <img
            src="/logo-gustoamate.jpg"
            alt="GustoaMate"
            className="logo-auth"
            onClick={() => navigate('/')}
          />
        </div>
        <div className="register-card">
          <h2 className="register-title">Sign Up</h2>
          <p className="register-subtitle">Crea tu cuenta para comenzar a comprar</p>

          {error && (
            <div className="error-alert">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="register-form">
            <div className="form-group">
              <label htmlFor="email">Email *</label>
              <input
                type="email"
                id="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                required
                className="form-input"
                placeholder="tu@email.com"
              />
            </div>

            <div className="form-group">
              <label htmlFor="password">Contraseña *</label>
              <input
                type="password"
                id="password"
                name="password"
                value={formData.password}
                onChange={handleChange}
                required
                className="form-input"
                placeholder="••••••••"
              />
            </div>

            <div className="form-group">
              <label htmlFor="first_name">Nombre *</label>
              <input
                type="text"
                id="first_name"
                name="first_name"
                value={formData.first_name}
                onChange={handleChange}
                required
                className="form-input"
                placeholder="Tu nombre"
              />
            </div>

            <div className="form-group">
              <label htmlFor="last_name">Apellido *</label>
              <input
                type="text"
                id="last_name"
                name="last_name"
                value={formData.last_name}
                onChange={handleChange}
                required
                className="form-input"
                placeholder="Tu apellido"
              />
            </div>

            <button
              type="submit"
              className="btn-submit"
              disabled={loading}
            >
              {loading ? 'Registrando...' : 'Registrarse'}
            </button>
          </form>

          <div className="register-footer">
            <p>¿Ya tienes una cuenta?</p>
            <button
              className="btn-signin"
              onClick={() => navigate('/login')}
            >
              Sign In
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default RegisterPage;
