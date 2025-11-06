import axios from 'axios';
import { getToken } from '../utils/auth';

const USERS_SERVICE_URL = process.env.REACT_APP_USERS_SERVICE_URL || 'http://localhost:8082';
const ITEMS_SERVICE_URL = process.env.REACT_APP_ITEMS_SERVICE_URL || 'http://localhost:8080';
const SEARCH_SERVICE_URL = process.env.REACT_APP_SEARCH_SERVICE_URL || 'http://localhost:8081';

// Crear instancias de axios para cada servicio
const usersAPI = axios.create({
  baseURL: USERS_SERVICE_URL,
});

const itemsAPI = axios.create({
  baseURL: ITEMS_SERVICE_URL,
});

const searchAPI = axios.create({
  baseURL: SEARCH_SERVICE_URL,
});

// Interceptor para agregar el token en todas las peticiones
const addAuthInterceptor = (apiInstance) => {
  apiInstance.interceptors.request.use(
    (config) => {
      const token = getToken();
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );
};

// Agregar interceptores a todas las instancias
addAuthInterceptor(usersAPI);
addAuthInterceptor(itemsAPI);
addAuthInterceptor(searchAPI);

export { usersAPI, itemsAPI, searchAPI };
