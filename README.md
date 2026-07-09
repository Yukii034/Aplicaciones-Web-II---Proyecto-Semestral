# Plataforma de Intercambios y Donaciones Locales
### API Backend · TDI-601 Aplicaciones Web II · 2026-1
### Universidad Laica Eloy Alfaro de Manabí

---

## ¿Qué es este proyecto?

API REST construida en Go que conecta a personas y negocios que tienen productos que ya no necesitan — ropa, electrodomésticos, alimentos próximos a vencer — con personas dispuestas a recibirlos o intercambiar algo a cambio.

El problema que resuelve: quien quiere donar o intercambiar algo no sabe quién lo necesita cerca, y quien necesita algo no sabe dónde buscarlo. Las soluciones actuales como Facebook Marketplace están orientadas a la compra-venta con dinero y dependen de contactos previos. Nuestra API ofrece un canal local y centralizado.

---

## Quién construyó cada módulo

| Módulo | Responsable | Entidades |
|---|---|---|
| Publicaciones e Inventario | Pierina Peñaherrera | `Inventario`, `Publicacion` |
| Reputación, Logros y Calificaciones | José Manuel Castillo | `Reputacion`, `Logro`, `Logro_Usuario`, `Calificacion` |
| Acuerdos, Items y Usuarios | Néstor Gallegos | `Acuerdo`, `AcuerdoItem`, `Usuario` |

---

## Stack Tecnológico

| Tecnología | Uso |
|---|---|
| Go 1.26+ | Lenguaje principal del backend |
| Chi Router | Manejo de rutas HTTP |
| GORM | ORM para acceso a datos |
| Golang-JWT + bcrypt | Autenticación y seguridad |
| Testify | Tests unitarios con mocks |
| Docker + docker-compose | Contenedores y orquestación |
| SQLite | Base de datos para desarrollo local |
| PostgreSQL | Base de datos para producción (Docker) |
| GitHub Actions | CI/CD (build → vet → test) |

---

## Cómo correrlo

### Requisitos
- Docker Desktop instalado y corriendo
- Git

### Levantar con Docker

```bash
git clone https://github.com/Yukii034/Aplicaciones-Web-II---Proyecto-Semestral.git
cd Aplicaciones-Web-II---Proyecto-Semestral
docker-compose up --build
```

La API quedará disponible en `http://localhost:8080`. Docker levanta automáticamente la API y PostgreSQL sin pasos adicionales.

Para cargar datos de ejemplo una vez que el contenedor esté corriendo:
```bash
curl -X POST http://localhost:8080/api/v1/seed
```

Para detener los contenedores:
```bash
docker-compose down
```

---

## Endpoints

### Auth (grupal)

| Método | Ruta | Descripción | Auth |
|---|---|---|---|
| POST | `/api/v1/auth/register` | Registrar usuario (tipo: persona, empresa, admin) | No |
| POST | `/api/v1/auth/login` | Login, devuelve JWT | No |
| POST | `/api/v1/seed` | Cargar datos de ejemplo | No |

---

### Módulo de Publicaciones e Inventario — Pierina Peñaherrerea

#### Inventario

| Método | Ruta | Descripción | Auth |
|---|---|---|---|
| GET | `/api/v1/inventario` | Listar todos los items | Bearer |
| GET | `/api/v1/inventario/{id}` | Obtener item por ID | Bearer |
| POST | `/api/v1/inventario` | Crear item en inventario | Bearer |
| PUT | `/api/v1/inventario/{id}` | Actualizar item | Bearer + Admin |
| DELETE | `/api/v1/inventario/{id}` | Borrar item | Bearer + Admin |

Ejemplo de body para crear:
```json
{
    "nombre": "Laptop Dell",
    "descripcion": "En buen estado",
    "categoria": "Tecnología",
    "estado_objeto": "usado",
    "disponibilidad": "disponible",
    "cantidad": 1
}
```

#### Publicaciones

| Método | Ruta | Descripción | Auth |
|---|---|---|---|
| GET | `/api/v1/publicaciones` | Listar publicaciones | Bearer |
| GET | `/api/v1/publicaciones/{id}` | Obtener publicacion por ID | Bearer |
| POST | `/api/v1/publicaciones` | Crear publicacion | Bearer |
| PUT | `/api/v1/publicaciones/{id}` | Actualizar publicacion | Bearer |
| DELETE | `/api/v1/publicaciones/{id}` | Borrar publicacion | Bearer |

Ejemplo de body para crear:
```json
{
    "titulo": "Cambio laptop por tablet",
    "tipo_oferta": "intercambio",
    "estado_publicacion": "disponible",
    "mensaje": "Acepto tablet de buena marca",
    "inventario_id": 1
}
```

> **Nota:** Al crear una publicacion se valida que el `inventario_id` exista. Si no existe devuelve 404.

---

### Módulo de Reputación, Logros y Calificaciones — José Manuel Castillo

| Método | Ruta | Descripción | Auth |
|---|---|---|---|
| GET | `/api/v1/reputaciones` | Listar reputaciones | Bearer |
| POST | `/api/v1/reputaciones` | Crear reputacion | Bearer |
| GET | `/api/v1/reputaciones/{id}` | Obtener por ID | Bearer |
| PUT | `/api/v1/reputaciones/{id}` | Actualizar | Bearer |
| DELETE | `/api/v1/reputaciones/{id}` | Borrar | Bearer |
| GET | `/api/v1/logros` | Listar logros | Bearer |
| POST | `/api/v1/logros` | Crear logro | Bearer |
| GET | `/api/v1/logros/{id}` | Obtener por ID | Bearer |
| PUT | `/api/v1/logros/{id}` | Actualizar | Bearer |
| DELETE | `/api/v1/logros/{id}` | Borrar | Bearer |
| GET | `/api/v1/logro_usuarios` | Listar logros de usuario | Bearer |
| POST | `/api/v1/logro_usuarios` | Asignar logro a usuario | Bearer |
| GET | `/api/v1/calificaciones` | Listar calificaciones | Bearer |
| POST | `/api/v1/calificaciones` | Crear calificacion | Bearer |

---

### Módulo de Acuerdos y Transacciones — Néstor Gallegos

| Método | Ruta | Descripción | Auth |
|---|---|---|---|
| GET | `/api/v1/acuerdos` | Listar acuerdos | Bearer |
| POST | `/api/v1/acuerdos` | Crear acuerdo | Bearer |
| GET | `/api/v1/acuerdos/{id}` | Obtener por ID | Bearer |
| PUT | `/api/v1/acuerdos/{id}` | Actualizar estado | Bearer |
| DELETE | `/api/v1/acuerdos/{id}` | Borrar | Bearer |
| GET | `/api/v1/acuerdo_items` | Listar items de acuerdo | Bearer |
| POST | `/api/v1/acuerdo_items` | Crear item de acuerdo | Bearer |
| GET | `/api/v1/usuarios` | Listar usuarios | Bearer |
| POST | `/api/v1/usuarios` | Crear usuario | Bearer |

---

## Arquitectura

El proyecto sigue una arquitectura en capas: cada request pasa por handler → service → repository antes de llegar a la base de datos. Los módulos de dominio están separados en sus propias carpetas dentro de `service/`.

```
Request HTTP
     ↓
  Handler          (internal/handlers/)
     ↓
  Service          (internal/service/modulo_*)
     ↓
  Repository       (internal/storage/)
     ↓
  Base de datos    (PostgreSQL via GORM)
```

```
Aplicaciones-Web-II---Proyecto-Semestral/
├── cmd/
│   └── api/
│       └── main.go              # punto de entrada, arma dependencias y rutas
│
├── internal/
│   ├── config/
│   │   └── config.go            # variables de entorno / .env
│   │
│   ├── handlers/                # capa HTTP: decodifica, llama al service, responde
│   │   ├── server.go            # struct Server + Deps
│   │   ├── respond.go           # RespondJSON/RespondError
│   │   ├── params.go
│   │   ├── auth.go
│   │   ├── seed.go
│   │   ├── acuerdo.go
│   │   ├── acuerdo_item.go
│   │   ├── usuario.go
│   │   ├── publicacion.go
│   │   ├── inventario.go
│   │   ├── reputacion.go
│   │   ├── logro.go
│   │   ├── logro_usuario.go
│   │   └── calificacion.go
│   │
│   ├── service/                 # lógica de negocio / validaciones
│   │   ├── auth.go              # JWT + bcrypt
│   │   ├── errores.go           # ErrNoEncontrado, ErrVacio, etc.
│   │   ├── modulo_aiu/          # Acuerdos, AcuerdoItem, Usuarios
│   │   │   ├── acuerdo.go
│   │   │   ├── acuerdo_item.go
│   │   │   └── usuario.go
│   │   ├── modulo_pi/           # Publicación, Inventario
│   │   │   ├── publicacion.go
│   │   │   └── inventario.go
│   │   └── modulo_rlc/          # Reputación, Logro, Logro_Usuario, Calificación
│   │       ├── reputacion.go
│   │       ├── logro.go
│   │       ├── logro_usuario.go
│   │       └── calificacion.go
│   │
│   ├── storage/                 # acceso a datos (GORM)
│   │   ├── almacen.go           # interfaces (Almacen, repositorios, etc.)
│   │   ├── almacen_sqlite.go
│   │   ├── factory.go           # abre sqlite o postgres + AutoMigrate
│   │   ├── seed.go              # datos de ejemplo
│   │   ├── acuerdo.go
│   │   ├── acuerdo_items.go
│   │   ├── usuario.go
│   │   ├── publicacion.go
│   │   ├── inventario.go
│   │   ├── reputacion.go
│   │   ├── logro.go
│   │   ├── logro_usuario.go
│   │   └── calificacion.go
│   │
│   ├── models/                  # structs compartidos (entidades)
│   │   ├── acuerdo.go
│   │   ├── acuerdo_item.go
│   │   ├── usuarios.go
│   │   ├── publicacion.go
│   │   ├── inventario.go
│   │   ├── reputacion.go
│   │   ├── logro.go
│   │   ├── logro_usuario.go
│   │   └── calificacion.go
│   │
│   ├── middleware/
│   │   ├── auth.go              # valida JWT, SoloAdmin (roles)
│   │   └── cors.go
│   │
│   └── httpserver/
│       └── httpserver.go        # wrapper de http.Server + graceful shutdown
│
├── postman/                     # colecciones de Postman
├── Dockerfile
├── docker-compose.yaml
├── .env.example
├── go.mod / go.sum
└── README.md
```

---

## Tests

```bash
# Correr todos los tests
go test ./... -cover

# Ver cobertura por módulo
go test ./internal/service/modulo_pi/... -cover
go test ./internal/storage/... -cover
go test ./internal/handlers/... -cover
```

Cobertura actual:
- `service/modulo_pi`: 95%+
- `service/modulo_rlc`: 95%+
- `service/modulo_aiu`: 50%+

---

## Roles y autenticación

La API usa JWT con dos roles:

- **persona / empresa:** puede crear inventario y publicaciones, ver todo
- **admin:** además puede actualizar y borrar inventario

Para usar rutas protegidas, incluye el token en el header:
```
Authorization: Bearer <token>
```

---

## Documentos del proyecto

Colección de Postman: `/postman/proyecto-semestral-coleccion.postman_collection.json`

Documentos adicionales (problema, entrevistas, diagrama ER, documento de cierre):
https://uleam-my.sharepoint.com/:f:/g/personal/e1350140990_live_uleam_edu_ec/IgDPcqsDRv5gRJwr6v39Be3IAdZLn4wA85WITu23YC422JA?e=QhxhUm