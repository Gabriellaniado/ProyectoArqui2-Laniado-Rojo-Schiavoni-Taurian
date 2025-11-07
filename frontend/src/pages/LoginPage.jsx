import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { userService } from '../services/userService';
import { setToken } from '../utils/auth';
import { saveCustomerID } from '../utils/auth';
import './LoginPage.css';

const LoginPage = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    email: '',
    password: ''
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
    if (!formData.email || !formData.password) {
      setError('Todos los campos son obligatorios');
      return;
    }

    try {
      setLoading(true);
      const response = await userService.login(formData);

      // Guardar el token en la cookie
      if (response.token) {
        setToken(response.token);
        saveCustomerID(response.customer_id);
        alert('Inicio de sesión exitoso');
        navigate('/');
      } else {
        setError('No se recibió el token de autenticación');
      }
    } catch (err) {
      console.error('Error during login:', err);
      setError('Email o contraseña incorrectos');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-page">
      <div className="login-container">
        <div className="login-header">
          <img
            src="/logo-gustoamate.jpg"
            alt="GustoaMate"
            className="logo-auth"
            onClick={() => navigate('/')}
          />
        </div>
        <div className="login-card">
          <h2 className="login-title">Sign In</h2>
          <p className="login-subtitle">Inicia sesión para continuar</p>

          {error && (
            <div className="error-alert">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="login-form">
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
                autoComplete="email"
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
                autoComplete="current-password"
              />
            </div>

            <button
              type="submit"
              className="btn-submit"
              disabled={loading}
            >
              {loading ? 'Iniciando sesión...' : 'Sign In'}
            </button>
          </form>

          <div className="login-footer">
            <p>¿No tienes una cuenta?</p>
            <button
              className="btn-signup"
              onClick={() => navigate('/registro')}
            >
              Sign Up
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
