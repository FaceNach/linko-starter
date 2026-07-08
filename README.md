# Linko

A toy URL shortener written in Go, built as the starter repository for **Logging and Telemetry**. It's intentionally small and a little messy, designed as a realistic playground for practicing observability: structured logging, metrics, distributed tracing, and profiling.

---

## 🇬🇧 English

### What is this project?

Linko is a minimal HTTP service that shortens long URLs into short codes and redirects visitors from the short code to the original destination. It ships with a tiny embedded HTML/JS frontend, HTTP Basic Auth-protected API endpoints, and a full observability stack (logs, metrics, traces, and profiles) wired in on purpose so it can be used as a teaching tool rather than a production-ready product.

Core features:
- **Shorten a URL** — `POST /api/shorten` generates a random 6-character short code and stores the mapping.
- **Redirect** — `GET /{shortCode}` looks up the code and redirects (HTTP 302) to the original URL, after verifying the destination is reachable.
- **List URLs** — `GET /api/urls` returns up to 10 stored links as JSON.
- **Stats** — `GET /api/stats` reports how many redirects happened and how many bytes were "saved" by shortening.
- **Login / auth** — Basic Auth against a small in-memory user table with bcrypt-hashed passwords.
- **Admin shutdown** — `POST /admin/shutdown` for graceful shutdown in non-production environments.

### Tech stack

**Language & runtime**
- [Go](https://go.dev/) 1.26 — the entire backend is plain Go, no web framework, using the standard library's `net/http.ServeMux` with Go 1.22+ method/path routing patterns.

**Storage**
- Flat-file storage (`internal/store`) — each shortened URL is a plain file on disk named after its short code, inside the `data/` directory. No external database.

**Authentication & security**
- `golang.org/x/crypto/bcrypt` — password hashing/verification for HTTP Basic Auth.
- `github.com/pkg/errors` — wrapped errors with stack traces for debugging.

**Observability — logging**
- `log/slog` (Go standard library) — structured logging throughout the app.
- `github.com/lmittmann/tint` — colorized, human-friendly log handler for local/terminal development.
- `github.com/mattn/go-isatty` — detects whether output is a TTY to decide on colored vs. plain/JSON logs.
- `gopkg.in/natefinch/lumberjack.v2` — log file rotation, compression, and retention.
- A custom middleware (`requestLogger` in `server.go`) logs every request with method, path, redacted client IP, duration, byte counts, user, and any error — with sensitive fields (passwords, keys, etc.) automatically redacted.

**Observability — metrics**
- `github.com/prometheus/client_golang` — exposes a `/metrics` endpoint (`http_requests_total` counter by method/path/status) scraped by Prometheus.

**Observability — tracing**
- [OpenTelemetry Go SDK](https://opentelemetry.io/) (`go.opentelemetry.io/otel`, `otel/sdk`, `otel/trace`) — distributed tracing instrumentation across handlers, auth, and outbound destination checks.
- `otlptracegrpc` — exports spans via OTLP/gRPC to a collector (Jaeger in this setup).
- `otelhttp` — automatic HTTP server instrumentation wrapping the whole mux.

**Observability — profiling**
- `net/http/pprof` — exposes CPU and memory profiling endpoints (protected by the same Basic Auth middleware); sample profiles (`cpu.prof`, `memory.prof`, `memory.svg`) are included in the repo for reference.

**Observability stack (via Docker Compose)**
- **Prometheus** — scrapes and stores metrics.
- **node-exporter** — host-level system metrics.
- **Grafana** — dashboards and visualization on top of Prometheus.
- **Jaeger** — trace collection and UI (OTLP receiver on 4317/4318, UI on 16686).

**Frontend**
- A single static, embedded (`//go:embed`) `index.html` — plain HTML/CSS/JS, no framework, styled as a retro green-on-black terminal.

**Testing**
- Go's built-in `testing` package (`requestLogger_test.go`).

### Project layout

```
main.go            - entrypoint, logger & tracing setup, graceful shutdown
server.go          - HTTP server, routing, request logging & metrics middleware
handlers.go        - HTTP handlers (index, login, shorten, redirect, list, stats)
auth.go            - Basic Auth middleware + bcrypt password validation
idMiddleware.go     - request ID middleware
destination.go     - validates that a shortened URL's target is reachable
tracing.go         - OpenTelemetry tracer provider setup
internal/store/    - flat-file backed URL storage
internal/build/    - build metadata (git SHA, build time)
internal/linkoerr/ - error wrapping with structured slog attributes
data/              - on-disk storage for shortened URLs
docker-compose.yaml - Prometheus, node-exporter, Grafana, Jaeger
prometheus.yml     - Prometheus scrape configuration
```

### Running it

```bash
# start the app
go run .

# start the observability stack
docker compose up
```

Then visit `http://localhost:8899`, Prometheus at `:9090`, Grafana at `:3000`, and Jaeger at `:16686`.

---

## 🇪🇸 Español

### ¿Qué es este proyecto?

Linko es un servicio HTTP mínimo que acorta URLs largas en códigos cortos y redirige a los visitantes desde ese código hacia el destino original. Incluye un pequeño frontend embebido en HTML/JS, endpoints de API protegidos con autenticación Basic Auth, y una stack completa de observabilidad (logs, métricas, trazas y profiling) integrada a propósito, ya que su objetivo es servir como herramienta educativa y no como un producto listo para producción.

Funcionalidades principales:
- **Acortar una URL** — `POST /api/shorten` genera un código corto aleatorio de 6 caracteres y guarda la relación.
- **Redirección** — `GET /{shortCode}` busca el código y redirige (HTTP 302) a la URL original, tras verificar que el destino esté disponible.
- **Listar URLs** — `GET /api/urls` devuelve hasta 10 enlaces guardados en formato JSON.
- **Estadísticas** — `GET /api/stats` reporta cuántas redirecciones ocurrieron y cuántos bytes se "ahorraron" al acortar.
- **Login / autenticación** — Basic Auth contra una pequeña tabla de usuarios en memoria con contraseñas hasheadas con bcrypt.
- **Apagado administrativo** — `POST /admin/shutdown` para un cierre controlado del servidor en entornos que no son de producción.

### Tecnologías utilizadas

**Lenguaje y runtime**
- [Go](https://go.dev/) 1.26 — todo el backend está escrito en Go puro, sin frameworks web, usando `net/http.ServeMux` de la librería estándar con los patrones de ruteo por método/path disponibles desde Go 1.22+.

**Almacenamiento**
- Almacenamiento en archivos planos (`internal/store`) — cada URL acortada se guarda como un archivo en disco nombrado según su código corto, dentro del directorio `data/`. No se usa una base de datos externa.

**Autenticación y seguridad**
- `golang.org/x/crypto/bcrypt` — hasheo y verificación de contraseñas para la autenticación Basic Auth.
- `github.com/pkg/errors` — errores envueltos con stack trace para facilitar el debugging.

**Observabilidad — logging**
- `log/slog` (librería estándar de Go) — logging estructurado en toda la aplicación.
- `github.com/lmittmann/tint` — handler de logs con colores, pensado para desarrollo local en terminal.
- `github.com/mattn/go-isatty` — detecta si la salida es una terminal (TTY) para decidir entre logs con color o en formato plano/JSON.
- `gopkg.in/natefinch/lumberjack.v2` — rotación, compresión y retención de archivos de log.
- Un middleware personalizado (`requestLogger` en `server.go`) registra cada request con método, path, IP del cliente (parcialmente redactada), duración, cantidad de bytes, usuario y errores — con campos sensibles (contraseñas, claves, etc.) redactados automáticamente.

**Observabilidad — métricas**
- `github.com/prometheus/client_golang` — expone un endpoint `/metrics` (contador `http_requests_total` por método/path/status) que es recolectado por Prometheus.

**Observabilidad — trazas (tracing)**
- [OpenTelemetry Go SDK](https://opentelemetry.io/) (`go.opentelemetry.io/otel`, `otel/sdk`, `otel/trace`) — instrumentación de trazas distribuidas en los handlers, la autenticación y las verificaciones de destino salientes.
- `otlptracegrpc` — exporta los spans vía OTLP/gRPC hacia un colector (Jaeger, en esta configuración).
- `otelhttp` — instrumentación automática del servidor HTTP que envuelve todo el mux.

**Observabilidad — profiling**
- `net/http/pprof` — expone endpoints de profiling de CPU y memoria (protegidos por el mismo middleware de Basic Auth); se incluyen perfiles de ejemplo (`cpu.prof`, `memory.prof`, `memory.svg`) en el repositorio como referencia.

**Stack de observabilidad (vía Docker Compose)**
- **Prometheus** — recolecta y almacena métricas.
- **node-exporter** — métricas del sistema a nivel de host.
- **Grafana** — dashboards y visualización sobre los datos de Prometheus.
- **Jaeger** — recolección de trazas y UI (receptor OTLP en 4317/4318, interfaz en 16686).

**Frontend**
- Un único `index.html` estático embebido (`//go:embed`) — HTML/CSS/JS puro, sin frameworks, con estilo retro de terminal verde sobre negro.

**Testing**
- Paquete `testing` incorporado de Go (`requestLogger_test.go`).

### Estructura del proyecto

```
main.go            - punto de entrada, configuración de logger y tracing, apagado ordenado
server.go          - servidor HTTP, ruteo, middlewares de logging y métricas
handlers.go        - handlers HTTP (index, login, shorten, redirect, list, stats)
auth.go            - middleware de Basic Auth + validación de contraseñas con bcrypt
idMiddleware.go     - middleware de ID de request
destination.go     - valida que el destino de una URL acortada esté disponible
tracing.go         - configuración del proveedor de trazas de OpenTelemetry
internal/store/    - almacenamiento de URLs basado en archivos planos
internal/build/    - metadata de build (SHA de git, hora de build)
internal/linkoerr/ - envoltura de errores con atributos estructurados para slog
data/              - almacenamiento en disco de las URLs acortadas
docker-compose.yaml - Prometheus, node-exporter, Grafana, Jaeger
prometheus.yml     - configuración de scraping de Prometheus
```

### Cómo ejecutarlo

```bash
# iniciar la aplicación
go run .

# iniciar la stack de observabilidad
docker compose up
```

Luego visitá `http://localhost:8899`, Prometheus en `:9090`, Grafana en `:3000`, y Jaeger en `:16686`.
