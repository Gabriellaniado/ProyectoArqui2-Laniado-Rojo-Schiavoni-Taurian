import { itemsAPI } from "./api";

export const salesService = {
  // Crear venta
  createSale: async (saleData) => {
    try {
      const response = await itemsAPI.post(
        "http://localhost:8080/sales",
        saleData
      );
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  // Obtener venta por ID
  getSaleById: async (saleId) => {
    try {
      const response = await itemsAPI.get(
        `http://localhost:8080/sales/${saleId}`
      );
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  // Obtener todas las ventas de un cliente
  getSalesByCustomerId: async (customerId) => {
    try {
      const response = await itemsAPI.get(
        `http://localhost:8080/sales/customer/${customerId}`
      );
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  // Actualizar venta
  updateSale: async (saleId, saleData) => {
    try {
      const response = await itemsAPI.put(
        `http://localhost:8080/sales/${saleId}`,
        saleData
      );
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },

  // Eliminar venta
  deleteSale: async (saleId) => {
    try {
      const response = await itemsAPI.delete(
        `http://localhost:8080/sales/${saleId}`
      );
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },
};
