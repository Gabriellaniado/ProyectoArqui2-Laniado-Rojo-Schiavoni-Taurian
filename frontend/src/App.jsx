import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { CartProvider} from "./context/CartContext";
import ProductsPage from './pages/ProductsPage';
import ProductDetailPage from './pages/ProductDetailPage';
import RegisterPage from './pages/RegisterPage';
import LoginPage from './pages/LoginPage';
import PurchasesPage from './pages/PurchasesPage';
import PurchaseDetailPage from './pages/PurchaseDetailPage';
import CartPage from "./pages/CartPage";
import CartDrawer from "./components/CartDrawer";
import AdminPage from './pages/AdminPage';
import NewProductPage from './pages/NewProductPage';
import EditProductPage from './pages/EditProductPage';
import './App.css';

function App() {
  return (
    <Router>
      <CartProvider>
        <div className="App">
          <Routes>
            <Route path="/" element={<ProductsPage />} />
            <Route path="/producto/:id" element={<ProductDetailPage />} />
            <Route path="/registro" element={<RegisterPage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/mis-compras" element={<PurchasesPage />} />
            <Route path="/compra/:id" element={<PurchaseDetailPage />} />
            <Route path="/carrito" element={<CartPage />} />
            <Route path="/admin" element={<AdminPage />} />
            <Route path="/admin/nuevo-producto" element={<NewProductPage />} />
            <Route path="/admin/editar-producto/:id" element={<EditProductPage />} />
          </Routes>
          <CartDrawer />
        </div>
      </CartProvider>
    </Router>
  );
}

export default App;
