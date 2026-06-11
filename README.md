# **Aplicaciones-Web-II---Proyecto-Semestral**

# **Descripcion:**
# El proyecto tratara de una API REST construida en GO que conecta a las personas y negocios que tienen productos que ya no necesitan como ropa, electrodomésticos, alimentos que estan próximos a vencer, etc. Con personas que si lo necesiten o estan dispuestas a intercambiar algo a cambio.

# **El problema que resuelve es simple:** quién quiere donar o intercambiar algo no sabe quien lo necesita cerca y quién necesita algo no sabe dónde poder buscarlo, las soluciones actuales como el Facebook Marketplace o grupos de WhatsApp no fueron diseñadAs principalmente con ese fin, están mas orientados a la compra-venta con dinero y depende de contactos previos o una cercanía fisica.

# **Nuestra API ofrece un canal local y centralizado con tres módulos de dominio:**
# **• Publicación:** Los usuarios registran productos disponibles para donaciones o intercambio, la app se encarga de marcarlos si esta disponible, reservado o entregado.

# **• Transacciones:** Cuando el cliente quiere un producto, se genera lo que es un acuerdos entre las dos partes que pasa por diferentes etapas: propuesto, aceptado y completado, ambos usuarios deben confirmar que la entrega ocurrió.

# **• Reputación:** Se construye un perfil de confianza de cada usuario que muestra sus puntos, niveles y logros desbloqueados, así la gente pueda confiar en la otra persona antes de poder hacer un trato y no genere desconfianza  

# **Tech Stack (Stack Tecnológico o Tecnologías a Utilizar):**
# **• Go 1.22+:** Lenguaje de programación de alto rendimiento utilizado para el desarrollo del backend, aprovechando su eficiencia en el manejo de concurrencia.
# **• Chi Router:** Framework ligero y minimalista para Go que permite gestionar las rutas de la API de forma rápida y con un control total sobre el flujo HTTP.
# **• GORM:** Herramienta de mapeo objeto-relacional (ORM) para Go que simplifica la interacción con la base de datos y la gestión de modelos complejos.
# **• Golang-JWT:** Librería encargada de la generación y validación de tokens de seguridad para garantizar una autenticación confiable en cada petición del usuario.
# **• Testify:** Suite de herramientas para pruebas unitarias en Go que facilita la escritura de tests descriptivos y asegura la estabilidad de la lógica de negocio.
# **• Docker:** Plataforma de contenedores que permite estandarizar el entorno de desarrollo y despliegue, asegurando que la API funcione igual en cualquier servidor.
# **• SQLite (Desarrollo):** Base de datos relacional ligera basada en archivos, ideal para realizar pruebas rápidas y desarrollo local sin dependencias externas.
# **• PostgreSQL (Producción):** Sistema de base de datos robusto y escalable utilizado para el entorno final, garantizando la persistencia y seguridad de la información.
# **• Git:** Sistema de control de versiones distribuido que permite registrar cada cambio realizado en el código y trabajar de forma segura en diferentes ramas.
# **• GitHub:** Plataforma en la nube que aloja el repositorio del proyecto, facilitando la colaboración entre los miembros del equipo y la integración del código.

# **URL a documentos anexos:**  https://uleam-my.sharepoint.com/:f:/g/personal/e1350140990_live_uleam_edu_ec/IgDPcqsDRv5gRJwr6v39Be3IAdZLn4wA85WITu23YC422JA?e=QhxhUm