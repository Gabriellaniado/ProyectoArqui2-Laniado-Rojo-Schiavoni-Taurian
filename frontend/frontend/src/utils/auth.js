import Cookies from 'js-cookie';
import { jwtDecode } from 'jwt-decode';

export const getToken = () => {
  return Cookies.get('token');
};

export const setToken = (token) => {
  Cookies.set('token', token, { expires: 7 }); // 7 días
};

export const removeToken = () => {
  Cookies.remove('token');
};

export const isAuthenticated = () => {
  const token = getToken();
  if (!token) return false;
  
  try {
    const decoded = jwtDecode(token);
    // Verificar si el token no ha expirado
    if (decoded.exp && decoded.exp * 1000 < Date.now()) {
      removeToken();
      return false;
    }
    return true;
  } catch (error) {
    removeToken();
    return false;
  }
};

export const getUserIdFromToken = () => {
  const token = getToken();
  if (!token) return null;
  
  try {
    const decoded = jwtDecode(token);
    return decoded.customer_id || decoded.user_id || decoded.id || decoded.sub;
  } catch (error) {
    return null;
  }
};

export const getDecodedToken = () => {
  const token = getToken();
  if (!token) return null;
  
  try {
    return jwtDecode(token);
  } catch (error) {
    return null;
  }
};
