# 🧉 E-Commerce de Mates – Proyecto de Microservicios

### Arquitectura de Software II – UCC

## 📖 Descripción general

Este proyecto consiste en el desarrollo de un **e-commerce de mates y accesorios** (mates, bombillas, yerberas, kits, etc.), implementado bajo una **arquitectura de microservicios**.  
Cada servicio tiene su propio dominio, base de datos y lógica independiente, comunicándose entre sí mediante **APIs REST** y **mensajería asíncrona (RabbitMQ)**.

El objetivo académico es **aplicar los conceptos vistos en clase**: separación por dominios, desacoplamiento, contenedorización con Docker, comunicación entre servicios y despliegue local completo.

---

## 👥 Integrantes

**Materia:** Arquitectura de Software II  
**Carrera:** Ingeniería en Sistemas de Información – Universidad Católica de Córdoba  
**Año:** 2025

**Alumnos:**

- Gabriel Laniado
- Candelaria Rojo
- Francisco Taurian
- Santino Schiavoni

**Docente:**

- Emiliano Kohmann

---

## 🧩 Arquitectura general

El sistema se compone de varios microservicios desarrollados en **Go (Gin)** y una infraestructura orquestada con **Docker Compose**.

### 🧱 Microservicios principales

| Servicio                | Descripción                                        | Base de datos |
| ----------------------- | -------------------------------------------------- | ------------- |
| `users-api`             | Manejo de usuarios, registro y autenticación (JWT) | MySQL         |
| `products-api`          | Catálogo de productos (mates, bombillas, etc.)     | MongoDB       |
| `search-api`            | Carrito y gestión de órdenes                       | SOLr          |
| `notifications-service` | Envío de emails y notificaciones push              | RabbitMQ      |

---

## 🧠 Objetivos de aprendizaje

- Aplicar los principios de **arquitectura orientada a microservicios**.
- Implementar **APIs RESTful** independientes y desacopladas.
- Integrar **Docker y Docker Compose** para levantar la infraestructura completa.
- Incorporar una **capa de mensajería (RabbitMQ)** para eventos entre servicios.
- Trabajar con **bases de datos heterogéneas** (relacional + no relacional).
- Diseñar un flujo básico de **autenticación, catálogo y compra**.

---
