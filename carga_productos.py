import requests
import json
import time

# --- Configuración ---
API_URL = "http://localhost:8080/items/"

# COLOCA AQUÍ TU TOKEN DE USUARIO
# (Asegúrate de copiarlo tal cual te lo da el login, sin espacios extra)
USER_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc19hZG1pbiI6dHJ1ZSwidXNlcl9pZCI6MSwiaXNzIjoiYmFja2VuZCIsInN1YiI6ImF1dGgiLCJleHAiOjE3NjQyNzcwOTYsIm5iZiI6MTc2NDI3NjQ5NiwiaWF0IjoxNzY0Mjc2NDk2LCJqdGkiOiIxIn0.-MYqG2uLUD2sdfBOJEMWerv1UU69vDlOiWFOrOgQn5U"

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
        "image_url": "https://cdn11.bigcommerce.com/s-3stx4pub31/images/stencil/1920w/products/1030/2867/rosamonte_premoum__68001.1648060634.png?c=2"
    },
    {
        "name": "Yerba Mate Playadito",
        "category": "Yerbas",
        "description": "Sabor suave y amigable, ideal para principiantes. Paquete de 1kg.",
        "price": 2950.00,
        "stock": 120,
        "image_url": "https://jumboargentina.vtexassets.com/arquivos/ids/711224/Yerba-Mate-Playadito-Suave-X1kg-1-854539.jpg?v=637938633804770000"
    },
    {
        "name": "Yerba Mate Taragüí",
        "category": "Yerbas",
        "description": "Energía y sabor clásico argentino. Elaborada con palo. Paquete de 1kg.",
        "price": 2750.00,
        "stock": 80,
        "image url": "https://th.bing.com/th/id/OIP.oB3t1bMWd00zl8uW2qlkAQHaHa?w=208&h=208&c=7&r=0&o=7&dpr=1.5&pid=1.7&rm=3"},
    {
        "name": "Yerba Mate Amanda Orgánica",
        "category": "Yerbas",
        "description": "Certificada orgánica, sin agroquímicos, sabor suave. Paquete de 500g.",
        "price": 3500.00,
        "stock": 50,
        "image_url": "https://th.bing.com/th/id/OIP.x-84SRRInKV3fdV8M20rYAHaHa?w=210&h=210&c=7&r=0&o=7&dpr=1.5&pid=1.7&rm=3"
    },
    {
        "name": "Yerba Mate Cruz de Malta",
        "category": "Yerbas",
        "description": "Elaborada con palo, molienda equilibrada y sabor tradicional. Paquete de 1kg.",
        "price": 2600.00,
        "stock": 75,
        "image_url": "https://th.bing.com/th/id/OIP.-QQZ0i2KEdBbzghjwFzqOAAAAA?w=152&h=203&c=7&r=0&o=7&dpr=1.5&pid=1.7&rm=3"
    },

    # Categoría: Mates
    {
        "name": "Mate de Calabaza",
        "category": "Mates",
        "description": "Clásico mate de calabaza curado, con virola de alpaca.",
        "price": 7500.00,
        "stock": 30,
        "image_url": "https://i.etsystatic.com/36685341/r/il/43d74c/4540646912/il_1080xN.4540646912_6cgs.jpg"
    },
    {
        "name": "Mate Imperial de Alpaca",
        "category": "Mates",
        "description": "Mate de calabaza forrado en cuero, con detalles y virola de alpaca.",
        "price": 15000.00,
        "stock": 15,
        "image_url": "https://th.bing.com/th/id/OIP.UPwQJWksMmPHYWIMp-B5bAHaGy?w=225&h=206&c=7&r=0&o=7&dpr=1.5&pid=1.7&rm=3"
    },
    {
        "name": "Mate de Madera (Palo Santo)",
        "category": "Mates",
        "description": "Mate de madera de palo santo, naturalmente aromático.",
        "price": 6000.00,
        "stock": 40,
        "image_url": "https://th.bing.com/th/id/OIP.w9Pey6Cuf8F-oXIR0veUEgHaDy?w=349&h=178&c=7&r=0&o=7&dpr=1.5&pid=1.7&rm=3"
    },
    {
        "name": "Mate de Silicona",
        "category": "Mates",
        "description": "Mate moderno, fácil de limpiar, no necesita curado. Varios colores.",
        "price": 4500.00,
        "stock": 60,
        "image_url": "https://th.bing.com/th/id/OIP.H8RS2uX3Dgdac5aG-t0tXQHaHM?w=217&h=211&c=7&r=0&o=7&dpr=1.5&pid=1.7&rm=3g"
    },
    {
        "name": "Mate de Acero Inoxidable",
        "category": "Mates",
        "description": "Doble capa térmica, mantiene la temperatura y es irrompible.",
        "price": 8000.00,
        "stock": 50,
        "image_url": "https://th.bing.com/th/id/OIP.pS_64JEGhQsEeFvYIzrehwHaHa?w=191&h=191&c=7&r=0&o=7&dpr=1.5&pid=1.7&rm=3"
    },

    # Categoría: Bombillas
    {
        "name": "Bombilla Pico de Loro",
        "category": "Bombillas",
        "description": "Clásica bombilla de alpaca, filtro de resorte. No se tapa.",
        "price": 3500.00,
        "stock": 80,
        "image_url": "https://th.bing.com/th/id/OIP.yQZ8c5ok-PMDjo49Dzw8bgHaHa?w=203&h=203&c=7&r=0&o=7&dpr=1.5&pid=1.7&rm=3"
    },
    {
        "name": "Bombilla Plana de Acero",
        "category": "Bombillas",
        "description": "Bombilla de acero inoxidable quirúrgico, fácil de limpiar.",
        "price": 4000.00,
        "stock": 100,
        "image_url":"https://th.bing.com/th/id/OIP.QPMT_yC49G-OY4HL9r4u2gHaHa?w=185&h=186&c=7&r=0&o=7&dpr=1.5&pid=1.7&rm=3"
    },
    {
        "name": "Bombilla con Filtro Desmontable",
        "category": "Bombillas",
        "description": "Permite una limpieza profunda del filtro. Incluye cepillo.",
        "price": 5000.00,
        "stock": 65,
        "image_url": "https://th.bing.com/th/id/OIP.1GpAibNAff_QaaZkFDWoLwHaJ4?w=146&h=195&c=7&r=0&o=7&dpr=1.5&pid=1.7&rm=3"
    },
    {
        "name": "Bombilla de Alpaca Cincelada",
        "category": "Bombillas",
        "description": "Diseño artesanal con detalles cincelados a mano.",
        "price": 9000.00,
        "stock": 25,
        "image_url": "https://th.bing.com/th/id/OIP.lVu51GMa7_X09s_eBIYCPwHaHa?w=215&h=215&c=7&r=0&o=7&dpr=1.5&pid=1.7&rm=3"
    },
    {
        "name": "Bombilla Cuchara",
        "category": "Bombillas",
        "description": "Filtro tipo cuchara, ideal para yerbas de molienda fina.",
        "price": 3800.00,
        "stock": 70,
        "image_url": "https://th.bing.com/th/id/OIP.IlWe3UNwxkwA3tmn147HXQHaJQ?w=158&h=198&c=7&r=0&o=7&dpr=1.5&pid=1.7&rm=3"
    },

    # Categoría: Accesorios
    {
        "name": "Termo Stanley 1L",
        "category": "Accesorios",
        "description": "Termo de acero inoxidable, mantiene el agua caliente 24hs. Pico cebador.",
        "price": 45000.00,
        "stock": 20,
        "image_url": "https://th.bing.com/th/id/OIP.8tpU8OWn8UOjspIBgqSF-QHaGo?w=238&h=212&c=7&r=0&o=7&dpr=1.5&pid=1.7&rm=3"
    },
    {
        "name": "Matera de Cuero",
        "category": "Accesorios",
        "description": "Bolso portatermo de cuero vacuno para llevar el set de mate completo.",
        "price": 22000.00,
        "stock": 30,
        "image_url": "https://th.bing.com/th/id/OIP.N-Hyk32C41op_bmyNj9dLAHaHa?w=202&h=202&c=7&r=0&o=7&dpr=1.5&pid=1.7&rm=3"
    },
    {
        "name": "Yerbero y Azucarera",
        "category": "Accesorios",
        "description": "Set de latas con pico vertedor para yerba y azúcar. Diseño pampa.",
        "price": 5500.00,
        "stock": 90,
        "image_url": "https://th.bing.com/th/id/OIP.12N8UygotjGg2_6im6J1-QHaHa?w=211&h=211&c=7&r=0&o=7&dpr=1.5&pid=1.7&rm=3"
    },
    {
        "name": "Pava Eléctrica ATMA",
        "category": "Accesorios",
        "description": "Pava eléctrica con corte automático a 85°C, temperatura ideal para mate.",
        "price": 28000.00,
        "stock": 40,
        "image_url": "https://th.bing.com/th/id/OIP.UYFgOYGccHou1IWwTgRETQHaIN?w=184&h=204&c=7&r=0&o=7&dpr=1.5&pid=1.7&rm=3"
    }
]

def cargar_items():
    """
    Función principal para iterar y enviar los datos a la API.
    """
    
    # Verificación simple para que no olvides poner el token
    if USER_TOKEN == "PEGAR_TU_TOKEN_LARGO_AQUI":
        print("⚠️ ERROR: No has configurado el USER_TOKEN en el script.")
        return

    print(f"--- Iniciando la carga de {len(items_para_cargar)} items a {API_URL} ---")

    # Configuración de los Headers con el Token
    # Se usa el estándar "Bearer". Si tu backend espera otro formato, modifícalo aquí.
    headers = {
        "Authorization": f"Bearer {USER_TOKEN}",
        "Content-Type": "application/json"
    }

    # Iterar sobre la lista y enviar cada item como POST
    for i, item in enumerate(items_para_cargar):
        try:
            # Realizar la solicitud POST pasando el json y los HEADERS
            response = requests.post(API_URL, json=item, headers=headers)

            # Verificar el código de estado de la respuesta
            if response.status_code in [200, 201]:
                print(f"({i+1}/{len(items_para_cargar)}) ÉXITO: Item '{item['name']}' cargado. (Status: {response.status_code})")
            elif response.status_code == 401:
                print(f"({i+1}/{len(items_para_cargar)}) ERROR DE AUTENTICACIÓN: Token inválido o expirado.")
                break # Detener si el token no sirve
            elif response.status_code == 403:
                print(f"({i+1}/{len(items_para_cargar)}) ERROR DE PERMISOS: El token es válido pero no tiene permisos.")
                break 
            else:
                # Si el servidor da otro error
                print(f"({i+1}/{len(items_para_cargar)}) ERROR al cargar '{item['name']}': Status {response.status_code}")
                print(f"   Respuesta del servidor: {response.text}")

        except requests.exceptions.ConnectionError:
            print(f"\n[ERROR DE CONEXIÓN] No se pudo conectar a {API_URL}.")
            print("Por favor, asegúrate de que tu servidor local (localhost:8080) esté corriendo.")
            break 
        except Exception as e:
            # Capturar cualquier otro error inesperado
            print(f"({i+1}/{len(items_para_cargar)}) Ocurrió un error inesperado al cargar '{item['name']}': {e}")

        # Opcional: una pequeña pausa
        # time.sleep(0.1)

    print("--- Carga de datos finalizada. ---")

if __name__ == "__main__":
    cargar_items()