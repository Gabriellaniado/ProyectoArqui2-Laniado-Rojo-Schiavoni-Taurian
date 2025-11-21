import React, { createContext, useContext, useState, useEffect } from 'react';
import { cartService } from '../services/cartService';
import { isAuthenticated, getCustomerId, getCustomerIDFromToken } from '../utils/auth';

const CartContext = createContext();

export const useCart = () => {
    const context = useContext(CartContext);
    if (!context) {
        throw new Error('useCart must be used within a CartProvider');
    }
    return context;
};

export const CartProvider = ({ children }) => {
    const [cart, setCart] = useState({
        items: [],
        total: 0,
        item_count: 0,
    });
    const [loading, setLoading] = useState(false);
    const [isOpen, setIsOpen] = useState(false);
    const [currentCustomerId, setCurrentCustomerId] = useState(null);

    // Solo cargar al montar si hay usuario autenticado
    useEffect(() => {
        const customerID = getCustomerIDFromToken();
        if (isAuthenticated() && customerID) {
            setCurrentCustomerId(customerID);
            loadCart();
        }
    }, []); // Solo se ejecuta una vez al montar

    // Función para inicializar el carrito desde el login
    // Cambiar la firma de initializeCart para recibir customerID
    const initializeCart = async (customerID, cartData) => {
        // Primero resetear el carrito anterior
        resetCart();

        // Actualizar el customerID
        setCurrentCustomerId(customerID);

        if (cartData && cartData.items && cartData.items.length > 0) {
            // Si viene data del backend, usarla
            setCart(cartData);
        } else {
            // Si no, cargar desde el servidor
            try {
                setLoading(true);
                const loadedCart = await cartService.getCart(customerID);
                setCart(loadedCart);
            } catch (error) {
                console.error('Error loading cart:', error);
                resetCart();
            } finally {
                setLoading(false);
            }
        }
    };

    // Función para resetear el carrito localmente
    const resetCart = () => {
        setCart({
            items: [],
            total: 0,
            item_count: 0,
        });
        setCurrentCustomerId(null);
        setIsOpen(false);
    };

    // Función para cargar el carrito desde el servidor
    const loadCart = async () => {
        try {
            setLoading(true);
            const customerID = getCustomerIDFromToken();

            if (!customerID) {
                resetCart();
                return;
            }

            const cartData = await cartService.getCart(customerID);
            setCart(cartData);
        } catch (error) {
            console.error('Error loading cart:', error);
            resetCart();
        } finally {
            setLoading(false);
        }
    };

    // Agregar item al carrito
    const addItem = async (itemID, quantity = 1) => {
        if (!isAuthenticated()) {
            alert('Debes iniciar sesión para agregar productos al carrito');
            return false;
        }

        try {
            setLoading(true);
            const customerID = getCustomerIDFromToken();
            const updatedCart = await cartService.addItem(customerID, itemID, quantity);
            setCart(updatedCart);
            return true;
        } catch (error) {
            console.error('Error adding item to cart:', error);
            const errorMessage = error.error || 'Error al agregar el producto al carrito';
            alert(errorMessage);
            return false;
        } finally {
            setLoading(false);
        }
    };

    // Actualizar cantidad de un item
    const updateItem = async (itemID, quantity) => {
        try {
            setLoading(true);
            const customerID = getCustomerIDFromToken();
            const updatedCart = await cartService.updateItem(customerID, itemID, quantity);
            setCart(updatedCart);
        } catch (error) {
            console.error('Error updating item:', error);
            const errorMessage = error.error || 'Error al actualizar el producto';
            alert(errorMessage);
        } finally {
            setLoading(false);
        }
    };

    // Eliminar item del carrito
    const removeItem = async (itemID) => {
        try {
            setLoading(true);
            const customerID = getCustomerIDFromToken();
            const updatedCart = await cartService.removeItem(customerID, itemID);
            setCart(updatedCart);
        } catch (error) {
            console.error('Error removing item:', error);
            alert('Error al eliminar el producto del carrito');
        } finally {
            setLoading(false);
        }
    };

    // Vaciar carrito
    const clearCart = async () => {
        try {
            setLoading(true);
            const customerID = getCustomerIDFromToken();
            if (customerID) {
                await cartService.clearCart(customerID);
            }
            resetCart();
        } catch (error) {
            console.error('Error clearing cart:', error);
            alert('Error al vaciar el carrito');
        } finally {
            setLoading(false);
        }
    };

    // Procesar checkout
    const checkout = async () => {
        try {
            setLoading(true);
            const customerID = getCustomerIDFromToken();
            const result = await cartService.checkout(customerID);
            resetCart();
            return result;
        } catch (error) {
            console.error('Error during checkout:', error);
            const errorMessage = error.error || 'Error al procesar la compra';
            throw new Error(errorMessage);
        } finally {
            setLoading(false);
        }
    };

    // Abrir/cerrar drawer del carrito
    const openCart = () => setIsOpen(true);
    const closeCart = () => setIsOpen(false);
    const toggleCart = () => setIsOpen(!isOpen);

    const value = {
        cart,
        loading,
        isOpen,
        initializeCart,
        resetCart,
        addItem,
        updateItem,
        removeItem,
        clearCart,
        checkout,
        loadCart,
        openCart,
        closeCart,
        toggleCart,
    };

    return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
};
