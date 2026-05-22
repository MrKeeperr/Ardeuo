import psycopg2
from psycopg2.extras import execute_values
from datetime import datetime, timedelta
import random
import math
import time

# Configuración de conexión
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
        nodo = random.choice(lista_nodos)
        hora_decimal = fecha_actual.hour + (fecha_actual.minute / 60.0)
        
        temp_base = 22 + 10 * math.sin(math.pi * (hora_decimal - 8) / 12)
        temperatura = round(temp_base + random.uniform(-1, 1), 2)
        hum_amb = round(max(30, min(95, 100 - (temperatura * 2) + random.uniform(-2, 2))), 2)
        hum_suelo = round(random.uniform(45.0, 65.0), 2)
        
        datos.append((fecha_actual, nodo, hum_suelo, temperatura, hum_amb))
        fecha_actual += timedelta(minutes=5)
        
    return datos, fecha_actual

def inyectar_millon_controlado():
    print("Conectando a la Bóveda de Datos...")
    conn = psycopg2.connect(**DB_CONFIG)
    cursor = conn.cursor()
    
    # --- OBTENER LOS NODOS DINÁMICAMENTE ---
    print("[🔍] Buscando nodos registrados en la base de datos...")
    cursor.execute("SELECT node_id FROM Nodos;")
    nodos_registrados = [fila[0] for fila in cursor.fetchall()]
    
    if not nodos_registrados:
        print("[❌] Error: No hay ningún nodo registrado en la tabla 'Nodos'.")
        return
        
    print(f"[✅] Se encontraron {len(nodos_registrados)} nodos para la simulación.")
    # ----------------------------------------
    
    META_REGISTROS = 1_000_000
    TAMANO_LOTE = 1000     
    TIEMPO_DESCANSO = 0.5  
    registros_insertados = 0
    
    fecha_simulada = datetime.now() - timedelta(days=1000) 
    
    query = """
        INSERT INTO telemetria (timestamp, node_id, humedad_suelo, temperatura_amb, humedad_amb)
        VALUES %s
    """
    
    print(f"\nIniciando inyección controlada de {META_REGISTROS:,} registros...")
    print(f"Configuración: Lotes de {TAMANO_LOTE} con {TIEMPO_DESCANSO}s de respiro.\n")
    
    try:
        tiempo_inicio = time.time()
        
        while registros_insertados < META_REGISTROS:
            lote, fecha_simulada = generar_lote_datos(fecha_simulada, TAMANO_LOTE, nodos_registrados)
            
            # Ejecución masiva en base de datos
            execute_values(cursor, query, lote)
            conn.commit()
            
            registros_insertados += TAMANO_LOTE
            porcentaje = (registros_insertados / META_REGISTROS) * 100
            
            # Print dinámico para rastrear el progreso en vivo
            print(f"\r[⚡] Lote insertado: {TAMANO_LOTE} | Progreso Total: {registros_insertados:,} / {META_REGISTROS:,} ({porcentaje:.2f}%)", end="", flush=True)
            
            # El descanso del procesador
            time.sleep(TIEMPO_DESCANSO)
            
        tiempo_total = round((time.time() - tiempo_inicio) / 60, 2)
        print(f"\n\n[✅] ¡Misión cumplida! 1,000,000 de registros inyectados de forma segura en {tiempo_total} minutos.")
            
    except Exception as e:
        print(f"\n\n[❌] Error durante la inyección: {e}")
        conn.rollback()
    finally:
        cursor.close()
        conn.close()

if __name__ == "__main__":
    inyectar_millon_controlado()