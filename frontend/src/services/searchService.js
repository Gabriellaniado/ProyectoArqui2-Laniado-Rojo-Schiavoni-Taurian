import { searchAPI } from "./api";

export const searchService = {
  // Buscar productos
  searchProducts: async (filters = {}) => {
    try {
      const params = new URLSearchParams();

      if (filters.query) params.append("query", filters.query);
      if (filters.id) params.append("id", filters.id);
      if (filters.name) params.append("name", filters.name);
      if (filters.minPrice !== undefined)
        params.append("minPrice", filters.minPrice);
      if (filters.maxPrice !== undefined)
        params.append("maxPrice", filters.maxPrice);
      if (filters.category) params.append("category", filters.category);
      if (filters.sort_by) params.append("sort_by", filters.sort_by);
      if (filters.page) params.append("page", filters.page);
      if (filters.count) params.append("count", filters.count);

      const response = await searchAPI.get(
        `http://localhost:8081/items?${params.toString()}`
      );
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },
};
