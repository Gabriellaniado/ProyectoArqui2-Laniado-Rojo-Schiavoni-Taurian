import requests
import json
import time

# --- Configuración ---
API_URL = "http://localhost:8080/items/"

# --- Datos de los 20 Productos ---
# 5 productos para 4 categorías (Yerbas, Mates, Bombillas, Accesorios)
items_para_cargar = [
    # Categoría: Yerbas
    {
        "name": "Yerba Mate Rosamonte",
        "category": "Yerbas",
        "description": "Yerba tradicional con palo, sabor intenso y duradero. Paquete de 1kg.",
        "price": 2800.00,
        "stock": 100,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },
    {
        "name": "Yerba Mate Playadito",
        "category": "Yerbas",
        "description": "Sabor suave y amigable, ideal para principiantes. Paquete de 1kg.",
        "price": 2950.00,
        "stock": 120,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },
    {
        "name": "Yerba Mate Taragüí",
        "category": "Yerbas",
        "description": "Energía y sabor clásico argentino. Elaborada con palo. Paquete de 1kg.",
        "price": 2750.00,
        "stock": 80,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },
    {
        "name": "Yerba Mate Amanda Orgánica",
        "category": "Yerbas",
        "description": "Certificada orgánica, sin agroquímicos, sabor suave. Paquete de 500g.",
        "price": 3500.00,
        "stock": 50,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },
    {
        "name": "Yerba Mate Cruz de Malta",
        "category": "Yerbas",
        "description": "Elaborada con palo, molienda equilibrada y sabor tradicional. Paquete de 1kg.",
        "price": 2600.00,
        "stock": 75,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },

    # Categoría: Mates
    {
        "name": "Mate de Calabaza",
        "category": "Mates",
        "description": "Clásico mate de calabaza curado, con virola de alpaca.",
        "price": 7500.00,
        "stock": 30,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },
    {
        "name": "Mate Imperial de Alpaca",
        "category": "Mates",
        "description": "Mate de calabaza forrado en cuero, con detalles y virola de alpaca.",
        "price": 15000.00,
        "stock": 15,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },
    {
        "name": "Mate de Madera (Palo Santo)",
        "category": "Mates",
        "description": "Mate de madera de palo santo, naturalmente aromático.",
        "price": 6000.00,
        "stock": 40,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },
    {
        "name": "Mate de Silicona",
        "category": "Mates",
        "description": "Mate moderno, fácil de limpiar, no necesita curado. Varios colores.",
        "price": 4500.00,
        "stock": 60,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },
    {
        "name": "Mate de Acero Inoxidable",
        "category": "Mates",
        "description": "Doble capa térmica, mantiene la temperatura y es irrompible.",
        "price": 8000.00,
        "stock": 50,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },

    # Categoría: Bombillas
    {
        "name": "Bombilla Pico de Loro",
        "category": "Bombillas",
        "description": "Clásica bombilla de alpaca, filtro de resorte. No se tapa.",
        "price": 3500.00,
        "stock": 80,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },
    {
        "name": "Bombilla Plana de Acero",
        "category": "Bombillas",
        "description": "Bombilla de acero inoxidable quirúrgico, fácil de limpiar.",
        "price": 4000.00,
        "stock": 100,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },
    {
        "name": "Bombilla con Filtro Desmontable",
        "category": "Bombillas",
        "description": "Permite una limpieza profunda del filtro. Incluye cepillo.",
        "price": 5000.00,
        "stock": 65,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },
    {
        "name": "Bombilla de Alpaca Cincelada",
        "category": "Bombillas",
        "description": "Diseño artesanal con detalles cincelados a mano.",
        "price": 9000.00,
        "stock": 25,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },
    {
        "name": "Bombilla Cuchara",
        "category": "Bombillas",
        "description": "Filtro tipo cuchara, ideal para yerbas de molienda fina.",
        "price": 3800.00,
        "stock": 70,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },

    # Categoría: Accesorios
    {
        "name": "Termo Stanley 1L",
        "category": "Accesorios",
        "description": "Termo de acero inoxidable, mantiene el agua caliente 24hs. Pico cebador.",
        "price": 45000.00,
        "stock": 20,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },
    {
        "name": "Kit de Limpieza (Ejemplo tuyo)",
        "category": "Accesorios",
        "description": "Kit completo de cepillos y herramientas para limpiar bombillas y mates.",
        "price": 3200.00,
        "stock": 55,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },
    {
        "name": "Matera de Cuero",
        "category": "Accesorios",
        "description": "Bolso portatermo de cuero vacuno para llevar el set de mate completo.",
        "price": 22000.00,
        "stock": 30,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },
    {
        "name": "Yerbero y Azucarera (Lata)",
        "category": "Accesorios",
        "description": "Set de latas con pico vertedor para yerba y azúcar. Diseño pampa.",
        "price": 5500.00,
        "stock": 90,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    },
    {
        "name": "Pava Eléctrica Matera",
        "category": "Accesorios",
        "description": "Pava eléctrica con corte automático a 85°C, temperatura ideal para mate.",
        "price": 28000.00,
        "stock": 40,
        "image_url": "https://inym.org.ar/imagenes/archivos/noticias/80133_imagen.jpg"
    }
]

def cargar_items():
    """
    Función principal para iterar y enviar los datos a la API.
    """
    print(f"--- Iniciando la carga de {len(items_para_cargar)} items a {API_URL} ---")

    # Iterar sobre la lista y enviar cada item como POST
    for i, item in enumerate(items_para_cargar):
        try:
            # Realizar la solicitud POST
            # 'requests' automáticamente serializa el dict a JSON
            # y establece el header 'Content-Type: application/json'
            response = requests.put(API_URL, json=item)

            # Verificar el código de estado de la respuesta
            # 201 (Created) es el código estándar para un POST exitoso que crea un recurso.
            # 200 (OK) también es común.
            if response.status_code in [200, 201]:
                print(f"({i+1}/{len(items_para_cargar)}) ÉXITO: Item '{item['name']}' cargado. (Status: {response.status_code})")
            else:
                # Si el servidor da un error (ej: 400, 500)
                print(f"({i+1}/{len(items_para_cargar)}) ERROR al cargar '{item['name']}': Status {response.status_code}")
                print(f"   Respuesta del servidor: {response.text}")

        except requests.exceptions.ConnectionError:
            print(f"\n[ERROR DE CONEXIÓN] No se pudo conectar a {API_URL}.")
            print("Por favor, asegúrate de que tu servidor local (localhost:8080) esté corriendo.")
            break # Detener el script si no hay conexión
        except Exception as e:
            # Capturar cualquier otro error inesperado
            print(f"({i+1}/{len(items_para_cargar)}) Ocurrió un error inesperado al cargar '{item['name']}': {e}")

        # Opcional: una pequeña pausa para no saturar el servidor
        # time.sleep(0.1)

    print("--- Carga de datos finalizada. ---")

if __name__ == "__main__":
    cargar_items()