import psycopg2
from psycopg2.extras import execute_values
from datetime import datetime, timedelta
import random
import math

# Configuración de conexión (Asegúrate de usar la contraseña real aquí también)
DB_CONFIG = {
    "host": "localhost",
    "port": 5432,
    "dbname": "ardeuo_db",
    "user": "admin_ardeuo",
    "password": "Neskuik@123"
}

def generar_lote_datos(fecha_inicio, cantidad, lista_nodos):
    """Genera datos usando la lista de nodos reales extraída de la BD."""
    datos = []
    fecha_actual = fecha_inicio
    
    for _ in range(cantidad):
        # Selecciona un nodo aleatorio de los que REALMENTE existen
        nodo = random.choice(lista_nodos)
        hora_decimal = fecha_actual.hour + (fecha_actual.minute / 60.0)
        
        temp_base = 22 + 10 * math.sin(math.pi * (hora_decimal - 8) / 12)
        temperatura = round(temp_base + random.uniform(-1, 1), 2)
        hum_amb = round(max(30, min(95, 100 - (temperatura * 2) + random.uniform(-2, 2))), 2)
        hum_suelo = round(random.uniform(45.0, 65.0), 2)
        
        datos.append((fecha_actual, nodo, hum_suelo, temperatura, hum_amb))
        fecha_actual += timedelta(minutes=5)
        
    return datos, fecha_actual

def ejecutar_prueba_piloto():
    print("[⚙️] Iniciando prueba piloto...")
    conn = None
    
    try:
        conn = psycopg2.connect(**DB_CONFIG)
        cursor = conn.cursor()
        
        print("[🔍] Buscando nodos registrados en la base de datos...")
        cursor.execute("SELECT node_id FROM Nodos;")
        nodos_registrados = [fila[0] for fila in cursor.fetchall()]
        
        if not nodos_registrados:
            print("[❌] Error: No hay ningún nodo registrado en la tabla 'Nodos'. Ejecuta primero el script SQL de creación.")
            return
            
        print(f"[✅] Se encontraron {len(nodos_registrados)} nodos: {nodos_registrados}")
        # --------------------------------------------------

        # 1. Generar exactamente 5 registros
        CANTIDAD_PRUEBA = 5
        fecha_simulada = datetime.now() - timedelta(hours=1) 
        lote, _ = generar_lote_datos(fecha_simulada, CANTIDAD_PRUEBA, nodos_registrados)
        
        query_insert = """
            INSERT INTO telemetria (timestamp, node_id, humedad_suelo, temperatura_amb, humedad_amb)
            VALUES %s
            RETURNING id, timestamp, node_id, humedad_suelo, temperatura_amb, humedad_amb;
        """
        
        print(f"\n[⚡] Inyectando {CANTIDAD_PRUEBA} registros en la base de datos...")
        query = """
            INSERT INTO telemetria (timestamp, node_id, humedad_suelo, temperatura_amb, humedad_amb)
            VALUES %s
        """
        execute_values(cursor, query, lote)
        conn.commit()
        print("[✅] Inserción exitosa.")
        
        # 3. Leer y mostrar los resultados
        print("\n[📊] Verificando los datos guardados en PostgreSQL:")
        cursor.execute("SELECT timestamp, node_id, humedad_suelo, temperatura_amb, humedad_amb FROM telemetria ORDER BY timestamp DESC LIMIT 5;")
        registros_guardados = cursor.fetchall()
        
        print("-" * 80)
        print(f"{'TIMESTAMP':<22} | {'MAC ADDRESS':<17} | {'HUM. SUELO (%)':<15} | {'TEMP (°C)':<10} | {'HUM. AMB (%)':<15}")
        print("-" * 80)
        
        for reg in registros_guardados:
            ts, nodo, h_suelo, temp, h_amb = reg
            ts_str = ts.strftime("%Y-%m-%d %H:%M:%S")
            print(f"{ts_str:<22} | {nodo:<17} | {h_suelo:<15} | {temp:<10} | {h_amb:<15}")
            
        print("-" * 80)
        print("\n[🚀] ¡Todo listo! Ya puedes aplicar esta misma lógica de extracción de nodos a tu inyector masivo.")
    except Exception as e:
        print(f"\n[❌] Error durante la prueba: {e}")
        if conn:
            conn.rollback()
    finally:
        if conn:
            cursor.close()
            conn.close()

if __name__ == "__main__":
    ejecutar_prueba_piloto()