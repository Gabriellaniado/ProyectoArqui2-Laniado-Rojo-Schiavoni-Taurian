import { usersAPI } from './api';

export const userService = {
  // Registro de usuario
  register: async (userData) => {
    try {
      const response = await usersAPI.post('http://localhost:8082/users', userData);
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  // Login
  login: async (credentials) => {
    try {
      const response = await usersAPI.post('http://localhost:8082/auth/login', credentials);
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  // Obtener usuario por ID
  getUserById: async (userId) => {
    try {
      const response = await usersAPI.get(`http://localhost:8082/users/${userId}`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  // Obtener usuario por email
  getUserByEmail: async (email) => {
    try {
      const response = await usersAPI.get(`http://localhost:8082/users/email/${email}`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  // Verificar token
  verifyToken: async () => {
    try {
      const response = await usersAPI.post('http://localhost:8082/auth/verify-token');
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  // Actualizar usuario
  updateUser: async (userId, userData) => {
    try {
      const response = await usersAPI.put(`http://localhost:8082/users/${userId}`, userData);
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  // Eliminar usuario
  deleteUser: async (userId) => {
    try {
      const response = await usersAPI.delete(`http://localhost:8082/users/${userId}`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  // Obtener todos los usuarios
  getUsers: async () => {
    try {
      const response = await usersAPI.get('http://localhost:8082/users');
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  }
};
