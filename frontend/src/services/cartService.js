import { itemsAPI } from './api';

export const cartService = {
    // Obtener carrito del usuario
    getCart: async (customerID) => {
        try {
            const response = await itemsAPI.get(`http://localhost:8080/cart/${customerID}`);
            return response.data;
        } catch (error) {
            // Si el carrito no existe, retornar carrito vacío
            if (error.response?.status === 404) {
                return {
                    customer_id: customerID,
                    items: [],
                    total: 0,
                    item_count: 0
                };
            }
            throw error.response?.data || error.message;
        }
    },

    // Agregar item al carrito
    addItem: async (customerID, itemID, quantity = 1) => {
        try {
            const response = await itemsAPI.post(
                `http://localhost:8080/cart/${customerID}/items`,
                {
                    item_id: itemID,
                    quantity: quantity
                }
            );
            return response.data;
        } catch (error) {
            throw error.response?.data || error.message;
        }
    },

    // Actualizar cantidad de un item
    updateItem: async (customerID, itemID, quantity) => {
        try {
            const response = await itemsAPI.put(
                `http://localhost:8080/cart/${customerID}/items/${itemID}`,
                {
                    quantity: quantity
                }
            );
            return response.data;
        } catch (error) {
            throw error.response?.data || error.message;
        }
    },

    // Eliminar item del carrito
    removeItem: async (customerID, itemID) => {
        try {
            const response = await itemsAPI.delete(
                `http://localhost:8080/cart/${customerID}/items/${itemID}`
            );
            return response.data;
        } catch (error) {
            throw error.response?.data || error.message;
        }
    },

    // Vaciar carrito
    clearCart: async (customerID) => {
        try {
            const response = await itemsAPI.delete(
                `http://localhost:8080/cart/${customerID}`
            );
            return response.data;
        } catch (error) {
            throw error.response?.data || error.message;
        }
    },

    // Procesar checkout
    checkout: async (customerID) => {
        try {
            const response = await itemsAPI.post(
                `http://localhost:8080/cart/${customerID}/checkout`
            );
            return response.data;
        } catch (error) {
            throw error.response?.data || error.message;
        }
    }
};