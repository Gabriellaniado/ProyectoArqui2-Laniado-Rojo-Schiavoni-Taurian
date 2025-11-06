import { itemsAPI } from './api';

export const productService = {
  // Obtener todos los productos (paginado)
  getAllProducts: async (page = 1, count = 9) => {
    try {
      const response = await itemsAPI.get(`http://localhost:8080/items?page=${page}&count=${count}`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  // Obtener producto por ID
  getProductById: async (productId) => {
    try {
      const response = await itemsAPI.get(`http://localhost:8080/items/${productId}`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  // Crear producto
  createProduct: async (productData) => {
    try {
      const response = await itemsAPI.post('http://localhost:8080/items', productData);
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  // Actualizar producto
  updateProduct: async (productId, productData) => {
    try {
      const response = await itemsAPI.put(`http://localhost:8080/items/${productId}`, productData);
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  // Eliminar producto
  deleteProduct: async (productId) => {
    try {
      const response = await itemsAPI.delete(`http://localhost:8080/items/${productId}`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  }
};
