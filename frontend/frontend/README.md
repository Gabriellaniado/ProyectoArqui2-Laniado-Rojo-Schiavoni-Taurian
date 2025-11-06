# E-commerce de Mates Argentinos - Frontend

Frontend en React para el e-commerce de mates argentinos con arquitectura de microservicios.

## 🛠️ Tecnologías Utilizadas

- React 18
- React Router DOM 6
- Axios
- js-cookie
- jwt-decode

## 📋 Prerequisitos

- Node.js (versión 14 o superior)
- npm o yarn
- Microservicios backend corriendo en:
  - Servicio de usuarios: `http://localhost:8082`
  - Servicio de productos: `http://localhost:8080`
  - Servicio de búsqueda: `http://localhost:8081`

## 🚀 Instalación

1. Instalar las dependencias:

```bash
npm install
```

2. Configurar las variables de entorno (opcional, ya están configuradas por defecto):

Editar el archivo `.env` si necesitas cambiar las URLs de los microservicios:

```
REACT_APP_USERS_SERVICE_URL=http://localhost:8082
REACT_APP_ITEMS_SERVICE_URL=http://localhost:8080
REACT_APP_SEARCH_SERVICE_URL=http://localhost:8081
```

## ▶️ Ejecución

Para iniciar el servidor de desarrollo:

```bash
npm start
```

La aplicación se abrirá en `http://localhost:3000`

## 📱 Funcionalidades

### Vistas Principales

1. **Vista de Productos** (`/`)
   - Lista de productos con imagen, nombre y precio
   - Barra de búsqueda
   - Paginación (3 productos por fila)
   - Botones de registro/login o mis compras (según autenticación)

2. **Vista de Producto Individual** (`/producto/:id`)
   - Información completa del producto
   - Selector de cantidad
   - Botón de compra con confirmación
   - Validación de autenticación

3. **Vista de Registro** (`/registro`)
   - Formulario con email, password, nombre y apellido
   - Validación de campos obligatorios
   - Redirección automática al login tras registro exitoso

4. **Vista de Login** (`/login`)
   - Formulario con email y password
   - Almacenamiento de token en cookie
   - Redirección a página principal tras login exitoso

5. **Vista de Mis Compras** (`/mis-compras`)
   - Lista de todas las compras del usuario
   - Información básica: número de compra, producto, fecha
   - Botón para ver detalle completo

6. **Vista de Detalle de Compra** (`/compra/:id`)
   - Información completa de la compra
   - Cantidad, precio total, fecha, IDs

## 🎨 Diseño

- Colores principales: 
  - Verde oscuro (botella): `#2d5016`, `#4a7c2a`
  - Crema/Blanco cálido: `#f5f3ed`, `#ffffff`
- Tipografía: Poppins (Google Fonts)
- Bordes redondeados
- Sombras suaves
- Animaciones de hover

## 🔐 Autenticación

- El token JWT se almacena en una cookie llamada `token`
- Validación automática de autenticación en rutas protegidas
- Redirección al login si no hay token válido
- Extracción del `customer_id` desde el token para operaciones

## 📂 Estructura del Proyecto

```
src/
├── components/          # Componentes reutilizables
│   ├── Header.jsx
│   ├── Header.css
│   ├── ProductCard.jsx
│   └── ProductCard.css
├── pages/              # Páginas/Vistas
│   ├── ProductsPage.jsx
│   ├── ProductsPage.css
│   ├── ProductDetailPage.jsx
│   ├── ProductDetailPage.css
│   ├── RegisterPage.jsx
│   ├── RegisterPage.css
│   ├── LoginPage.jsx
│   ├── LoginPage.css
│   ├── PurchasesPage.jsx
│   ├── PurchasesPage.css
│   ├── PurchaseDetailPage.jsx
│   └── PurchaseDetailPage.css
├── services/           # Servicios HTTP
│   ├── api.js
│   ├── userService.js
│   ├── productService.js
│   ├── salesService.js
│   └── searchService.js
├── utils/              # Utilidades
│   └── auth.js
├── App.jsx
├── App.css
├── index.js
└── index.css
```

## 🌐 Endpoints Utilizados

### Usuarios (Puerto 8082)
- `POST /users` - Registro de usuario
- `POST /auth/login` - Login
- `GET /users/:id` - Obtener usuario
- `GET /users/email/:email` - Obtener usuario por email

### Productos (Puerto 8080)
- `GET /items/:id` - Obtener producto por ID
- `POST /sales` - Crear venta
- `GET /sales/customer/:customerID` - Obtener ventas del cliente

### Búsqueda (Puerto 8081)
- `GET /search?query=...&page=...&count=...` - Buscar productos

## 🔧 Scripts Disponibles

- `npm start` - Inicia el servidor de desarrollo
- `npm run build` - Construye la aplicación para producción
- `npm test` - Ejecuta los tests
- `npm run eject` - Expone la configuración de Create React App

## 📝 Notas Importantes

1. Todos los endpoints requieren el token de autenticación en el header (excepto registro y login)
2. Las respuestas del backend vienen en formato `{ item: {...} }`
3. El token debe parsearse para obtener el `customer_id`
4. Los productos se muestran de 9 en 9 (3x3) por página
5. Las validaciones de compra verifican la autenticación antes de procesar

## 🐛 Troubleshooting

- Si hay errores de CORS, asegúrate de que los microservices tengan CORS habilitado
- Si el token no se guarda, verifica que las cookies estén habilitadas en el navegador
- Si no se cargan productos, verifica que el servicio de búsqueda esté corriendo en el puerto 8081

## 📄 Licencia

Este proyecto es parte de un sistema de e-commerce de mates argentinos.
