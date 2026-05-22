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

Fase C: Conexión con svelte (Frontend)
Configurar CORS: Activar el middleware go-chi/cors para evitar que el navegador bloquee las peticiones cuando tu frontend intente conectarse.

Endpoints del Dashboard: Crear las rutas definitivas que tu panel de control necesita (ej. GET /api/v1/telemetry, GET /api/v1/crops, POST /api/v1/irrigation).

--------------------------------------------------
PARTE 2

ADMIN SISTEMA CARACTERISTICAS
1. Gestión de "Tenants" (Multitenancy)
Esta es la función más importante. Dado que tu arquitectura está diseñada para escalar, el SuperAdmin necesita ver el sistema como un conjunto de clientes:

Directorio de Propietarios (Tenants): Ver cuántas fincas/propietarios están activos en el sistema.

Provisionamiento: Capacidad de activar o suspender el acceso a nuevos clientes.

Cuotas de Uso: Monitorear qué finca está generando más datos (qué finca está "gastando" más recursos de tu base de datos TimescaleDB) para ajustar planes de precios.

2. Salud de la Infraestructura (Observabilidad)
El SuperAdmin no mira humedad o temperatura, mira si el sistema está vivo:

Estado de los Servicios: Ver si los contenedores (Go, Python, Base de datos, MQTT broker) están arriba o abajo.

Logs Centralizados: Un buscador de errores de toda la plataforma. Si el backend de Go lanza un error 500, el SuperAdmin debe verlo aquí sin tener que entrar a los logs de los contenedores por terminal.

Salud del Almacenamiento: Monitorear el crecimiento de la base de datos (como ya estamos haciendo con la consulta de bytes) para anticipar cuándo necesitarás más disco duro.

3. Analytics del Ecosistema (Business Intelligence)
Adopción Tecnológica: ¿Cuántos sensores hay conectados en total a través de todas las fincas?

Análisis de Rendimiento: ¿Qué finca o qué tipo de cultivo está usando mejor la tecnología? (esto es oro para el marketing del Proyecto Demeter).

---------------------------

PROPIETARIO CARACTERISTICAS
1. Centro de Control Administrativo (Gestión)
Como el administrador es quien rige la operación, necesita controlar quién y qué está dentro del sistema:

Gestión de Usuarios y Roles: Capacidad de crear cuentas para operarios, asignar niveles de acceso y auditar quién ha realizado qué acciones (usando tu tabla responsables_eventos).

Inventario y Configuración IoT: Registro y edición de los nodos ESP32, asignación de sensores a cultivos específicos y configuración de los umbrales de alerta (ej. "si la humedad baja del 30%, disparar alerta").

Gestión de Cultivos: Creación, edición y cierre de ciclos de cultivo. Aquí es donde se define la lógica de negocio que alimenta tu modelo relacional en PostgreSQL.

2. Tablero de Inteligencia Agrícola (Big Data Analytics)
Utilizando la capacidad de TimescaleDB, el administrador debe ver tendencias, no solo valores instantáneos:

Análisis de Rendimiento Hídrico: Comparativa entre el consumo de agua real (eventos de riego) vs. los requerimientos hídricos del cultivo.

Predicción y Estabilidad: Visualización de promedios móviles de temperatura y humedad para identificar microclimas o zonas de estrés hídrico que el operario no percibe a pie de campo.

Exportación de Reportes: Generación de resúmenes mensuales (en PDF o CSV) para auditorías, presupuestos o toma de decisiones sobre insumos.

3. Centro de Trazabilidad y Auditoría (Audit Log)
Esta es la feature que "cierra el círculo" de tu arquitectura:

Historial de Decisiones: Un registro completo de qué evento ocurrió, cuándo y quién fue el responsable. Esto responde a la pregunta clave: "¿Por qué se activó el riego a las 3:00 AM y quién lo autorizó?".

Trazabilidad de Nodos: Si un sensor empieza a dar lecturas erróneas, el admin debe poder ver el historial de errores del nodo y decidir si necesita mantenimiento o reemplazo.