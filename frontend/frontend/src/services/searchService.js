import { searchAPI } from './api';

export const searchService = {
  // Buscar productos
  searchProducts: async (filters = {}) => {
    try {
      const params = new URLSearchParams();

      if (filters.query) params.append('query', filters.query);
      if (filters.id) params.append('id', filters.id);
      if (filters.name) params.append('name', filters.name);
      if (filters.min_price !== undefined) params.append('min_price', filters.min_price);
      if (filters.max_price !== undefined) params.append('max_price', filters.max_price);
      if (filters.category) params.append('category', filters.category);
      if (filters.sort_by) params.append('sort_by', filters.sort_by);
      if (filters.page) params.append('page', filters.page);
      if (filters.count) params.append('count', filters.count);

      const response = await searchAPI.get(`http://localhost:8081/items?${params.toString()}`);
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  }
};
