# Matrix QR Decomposition App

Aplicación web para descomposición QR de matrices utilizando el proceso de **Gram-Schmidt**. Desplegada en **AWS EKS Auto Mode**.

## Demo

**URL:** https://test-interseguro.duckdns.org/home

## Arquitectura

| Servicio | Puerto | Tecnología | Descripción |
|----------|--------|------------|-------------|
| `front` | 3000 | React + TanStack SSR + Vite | UI con formulario de matrices y visualización de resultados |
| `nginx` | 8080 | NGINX Alpine | Sidecar: proxy `/api` → Go API, `/` → Front |
| `go-api` | 3000 | Go + Fiber + GORM | API principal: auth JWT, QR decomposition, persistencia |
| `nodejs-api` | 3001 | Express + TypeScript | Microservicio de estadísticas matriciales |
| `postgres` | 5432 | PostgreSQL 16 | Base de datos con schema `finance` |

### Flujo de descomposición QR

1. El usuario ingresa una matriz `m×n` (m ≥ n) en el frontend
2. `POST /api/matrix/qr` (autenticado) envía la matriz al Go API
3. Go API ejecuta Gram-Schmidt: calcula `Q` (ortogonal) y `R` (triangular superior)
4. Opcionalmente consulta al Node.js API para estadísticas adicionales
5. Resultado se persiste en PostgreSQL y se devuelve al frontend

### Endpoints de la API

| Método | Ruta | Auth | Descripción |
|--------|------|------|-------------|
| POST | `/api/auth/register` | - | Registrar usuario |
| POST | `/api/auth/login` | - | Iniciar sesión (devuelve JWT) |
| POST | `/api/auth/logout` | JWT | Cerrar sesión |
| POST | `/api/auth/refresh-token` | - | Refrescar JWT |
| GET | `/api/users/:id` | JWT | Obtener usuario |
| PATCH | `/api/users/:id` | JWT | Actualizar usuario |
| DELETE | `/api/users/:id` | JWT | Eliminar usuario |
| POST | `/api/matrix/qr` | JWT | Calcular descomposición QR |
| GET | `/api/matrix/qr` | JWT | Listar resultados |
| GET | `/api/matrix/qr/:id` | JWT | Obtener resultado por ID |

## Ejecución local

### Prerrequisitos

- Docker + Docker Compose

### Pasos

```bash
# 1. Copiar variables de entorno
cp .env.example .env

# 2. Iniciar todos los servicios
docker compose up -d

# 3. Acceder
curl http://localhost/

# 4. Detener
docker compose down -v
```

## Despliegue en Kubernetes (Helm)

```bash
# Instalar/actualizar
helm upgrade --install test-interseguro ./helm -n test-interseguro --create-namespace

# Ver pods
kubectl get pods -n test-interseguro

# Ver logs
kubectl logs -n test-interseguro deployment/go-api
kubectl logs -n test-interseguro deployment/nodejs-api
kubectl logs -n test-interseguro deployment/front
kubectl logs -n test-interseguro statefulset/postgres
```

## Estructura del proyecto

```
├── front/              # React + TanStack SSR
├── go-api/             # Go Fiber API
├── nodejs-api/         # Express TypeScript API
├── nginx/              # Config NGINX
├── postgres/           # Init scripts SQL
├── helm/               # Chart Helm
├── k8s/                # Manifiestos K8s raw
├── docker-compose.yml
└── .env.example
```

## Decisiones técnicas

- **EKS Auto Mode:** Los nodos están en subnets privadas. Se requieren VPC Endpoints (ECR, EC2, STS, ELB) y un Gateway endpoint para S3.
- **EBS CSI Driver:** No compatible con EKS Auto (node affinity bloquea nodos `compute-type=auto`). Se usa `emptyDir` para PostgreSQL en dev.
- **NGINX como sidecar:** En K8s se despliega un contenedor NGINX junto al frontend para enrutar `/api` al Go API, evitando un Ingress Controller.
