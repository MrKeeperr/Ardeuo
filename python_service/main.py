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

db_connection = conectar_db()

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
    try:
        payload = msg.payload.decode('utf-8')
        datos = json.loads(payload)
        
        node_id = datos.get("node_id")
        humedad_suelo = datos.get("humedad_suelo")
        temperatura_amb = datos.get("temperatura_amb")
        humedad_amb = datos.get("humedad_amb")
        timestamp = datetime.now()

        # Inserción a la base de datos
        if db_connection:
            cursor = db_connection.cursor()
            query = """
                INSERT INTO Telemetria (node_id, timestamp, humedad_suelo, temperatura_amb, humedad_amb)
                VALUES (%s, %s, %s, %s, %s)
            """
            valores = (node_id, timestamp, humedad_suelo, temperatura_amb, humedad_amb)
            
            cursor.execute(query, valores)
            db_connection.commit()
            cursor.close()
            
            print(f"💾 [GUARDADO] Nodo: {node_id} | Humedad Suelo: {humedad_suelo}% | Temp: {temperatura_amb}°C | Timestamp: {timestamp}")

    except json.JSONDecodeError:
        print("⚠️ [ERROR] El mensaje recibido no es un JSON válido.")
    except Exception as e:
        print(f"❌ [ERROR] Procesando el mensaje: {e}")
        if db_connection:
            db_connection.rollback()


if __name__ == "__main__":
    cliente_mqtt = mqtt.Client()
    cliente_mqtt.on_connect = on_connect # inicializamos la conexión al broker MQTT
    cliente_mqtt.on_message = on_message # inicializamos la función que se ejecutará al recibir un mensaje MQTT

    print(f"🚀 Iniciando conexión python-mqtt-iot...")
    try:
        cliente_mqtt.connect(MQTT_BROKER, MQTT_PORT, 60)
        cliente_mqtt.loop_forever()
    except KeyboardInterrupt:
        print("\n🛑 Apagando el servicio...")
    finally:
        if db_connection:
            db_connection.close()
            print("🔒 Conexión a DB cerrada.")