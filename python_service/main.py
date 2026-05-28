import os
import json
import psycopg2
import paho.mqtt.client as mqtt
from datetime import datetime
from dotenv import load_dotenv

load_dotenv()

DB_HOST = os.getenv("DB_HOST")
DB_PORT = os.getenv("DB_PORT")
DB_NAME = os.getenv("DB_NAME")
DB_USER = os.getenv("DB_USER")
DB_PASS = os.getenv("DB_PASS")

MQTT_BROKER = os.getenv("MQTT_BROKER")
MQTT_PORT = int(os.getenv("MQTT_PORT", 1883))
MQTT_TOPIC = os.getenv("MQTT_TOPIC")
MQTT_USER = os.getenv("MQTT_USER")
MQTT_PASS = os.getenv("MQTT_PASS")

def conectar_db():
    """Función para conectar a la base de datos TimescaleDB"""
    try:
        conn = psycopg2.connect(
            host=DB_HOST,
            port=DB_PORT,
            dbname=DB_NAME,
            user=DB_USER,
            password=DB_PASS
        )
        print("✅ [DB] Conectado exitosamente a TimescaleDB")
        return conn
    except Exception as e:
        print(f"❌ [DB] Error al conectar a la base de datos: {e}")
        return None

def verificar_dispositivo(conn, node_id):
    """Verifica si el dispositivo (node_id) existe en la tabla Dispositivos"""
    try:
        cursor = conn.cursor()
        cursor.execute("SELECT node_id FROM nodos WHERE node_id = %s", (node_id,))
        existe = cursor.fetchone() is not None
        cursor.close()
        return existe
    except Exception as e:
        print(f"❌ [DB] Error verificando dispositivo: {e}")
        return False

def on_connect(client, userdata, flags, rc):
    """
    Función que se ejecuta cuando se conecta al broker MQTT
    """
    if rc == 0:
        print(f"✅ [MQTT] Conectado a Mosquitto en {MQTT_BROKER}")
        client.subscribe(MQTT_TOPIC)
        print(f"📡 [MQTT] Escuchando mensajes en el topic: {MQTT_TOPIC}")
    else:
        print(f"❌ [MQTT] Error de conexión. Código: {rc}")

def on_message(client, userdata, msg):
    """
    Función que se ejecuta cuando se recibe un mensaje en el topic MQTT
    toma los datos del mensaje, y los carga a la BD
    """
    conn = userdata
    try:
        payload = msg.payload.decode('utf-8')
        datos = json.loads(payload)
        
        node_id = datos.get("node_id")
        humedad_suelo = datos.get("hum_suelo")
        temperatura_amb = datos.get("temp")
        humedad_amb = datos.get("hum_amb")
        timestamp = datetime.now()

        if not node_id:
            print("⚠️ [ERROR] El mensaje no contiene node_id.")
            return

        # Verificar que el dispositivo existe antes de insertar
        if conn is None:
            print("❌ [DB] No hay conexión a la base de datos.")
            return

        if not verificar_dispositivo(conn, node_id):
            print(f"⚠️ [DISPOSITIVO] Nodo {node_id} no registrado. Mensaje descartado.")
            return

        # Inserción a la base de datos
        cursor = conn.cursor()
        query = """
            INSERT INTO Telemetria (node_id, timestamp, humedad_suelo, temperatura_amb, humedad_amb)
            VALUES (%s, %s, %s, %s, %s)
        """
        valores = (node_id, timestamp, humedad_suelo, temperatura_amb, humedad_amb)
        
        cursor.execute(query, valores)
        conn.commit()
        cursor.close()
        
        print(f"💾 [GUARDADO] Nodo: {node_id} | Humedad Suelo: {humedad_suelo}% | Temp: {temperatura_amb}°C | Timestamp: {timestamp}")

    except json.JSONDecodeError:
        print("⚠️ [ERROR] El mensaje recibido no es un JSON válido.")
    except Exception as e:
        print(f"❌ [ERROR] Procesando el mensaje: {e}")
        if conn:
            conn.rollback()


if __name__ == "__main__":
    db_connection = conectar_db()

    cliente_mqtt = mqtt.Client()
    cliente_mqtt.username_pw_set(MQTT_USER, MQTT_PASS)
    cliente_mqtt.on_connect = on_connect
    cliente_mqtt.on_message = on_message

    print("🚀 Iniciando conexión python-mqtt-iot...")
    try:
        cliente_mqtt.connect(MQTT_BROKER, MQTT_PORT, 60)
        cliente_mqtt.user_data_set(db_connection)
        cliente_mqtt.loop_forever()
    except KeyboardInterrupt:
        print("\n🛑 Apagando el servicio...")
    finally:
        if db_connection:
            db_connection.close()
            print("🔒 Conexión a DB cerrada.")