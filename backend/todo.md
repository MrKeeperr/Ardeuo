El Proyecto Demeter es una iniciativa de agricultura de precisión que busca transformar la gestión agrícola empírica en una operación basada en datos de alta precisión. Ante la escasez hídrica y los retos climáticos en la región de Sabana Centro, el sistema actúa como una herramienta de toma de decisiones inteligente.

¿Cuál es la visión del proyecto?
El eje central es la optimización de recursos hídricos. En lugar de regar por intuición o por horarios fijos, el sistema utiliza datos reales capturados por nodos IoT (ESP32 con sensores de humedad FC-28 y temperatura DHT22) para regar exactamente cuando y donde el cultivo lo necesita. Esto no solo maximiza la producción, sino que democratiza la tecnología para microempresas rurales, garantizando sostenibilidad en un sector bajo presión demográfica.

¿Qué estamos haciendo actualmente?
Estamos construyendo el "cerebro" y el "sistema nervioso" de esta plataforma. Nuestra arquitectura está migrando hacia un ecosistema robusto, escalable y aislado utilizando Docker para la infraestructura y un enfoque de Big Data para el manejo de la telemetría.

Aquí está el desglose de lo que hemos avanzado técnicamente:

Infraestructura de Datos (PostgreSQL + TimescaleDB):

Hemos implementado una arquitectura híbrida donde PostgreSQL gestiona los metadatos relacionales (propietarios, cultivos, roles) y TimescaleDB maneja de forma optimizada los millones de registros de telemetría (series temporales) que generan tus sensores.

Desarrollo del Backend en Go (Golang):

Estamos construyendo un microservicio de backend de alto rendimiento en Go, siguiendo estándares de arquitectura empresarial (layout cmd/ e internal/).

Ya logramos configurar la conexión con la base de datos utilizando un pool de conexiones (pgxpool), garantizando que la API sea rápida y capaz de manejar concurrencia.

Hemos configurado un servidor web profesional con el enrutador Chi, aplicando middlewares de seguridad y logs para asegurar la trazabilidad.

Seguridad y Autenticación:

Estamos implementando una arquitectura de seguridad basada en JWT (JSON Web Tokens) y encriptación Bcrypt para proteger las contraseñas, permitiendo diferenciar roles (Propietarios vs. Operarios) sin saturar la base de datos con usuarios de sistema innecesarios.

Fase A: Seguridad y Autenticación (Lo más urgente)
Lógica de Encriptación: Crear una función de utilidad que use bcrypt para encriptar contraseñas al crear un usuario, y compararlas al hacer login.

Generador de JWT: Crear la función que firme y entregue los "pases VIP" (Tokens JWT) cuando las credenciales sean correctas.

Middleware de Protección: Escribir el guardián de Chi que intercepte el r.Route privado y bloquee a cualquiera que no envíe un token válido.

Fase B: La Capa de Base de Datos (Repository)
Dejar de usar datos "falsos" o de prueba.

Escribir las funciones que ejecutan código SQL real (SELECT, INSERT) a través de pgxpool para buscar usuarios, validar logins y consultar la telemetría de TimescaleDB.

Fase C: Conexión con React (Frontend)
Configurar CORS: Activar el middleware go-chi/cors para evitar que el navegador bloquee las peticiones cuando tu frontend intente conectarse.

Endpoints del Dashboard: Crear las rutas definitivas que tu panel de control necesita (ej. GET /api/v1/telemetry, GET /api/v1/crops, POST /api/v1/irrigation).