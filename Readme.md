# 📬 ose-postman

A gRPC-based Email Template and Messaging Service built with CQRS, PostgreSQL, MongoDB, and NATS.

## ✨ Features

- CRUD for templates (HTML + plain text)
- CRUD from email (also send mail via smtp)
- Send dynamic emails using templates and variables
- NATS pub/sub for async email dispatching
- CQRS pattern with:
  - **PostgreSQL** for command/write operations
  - **MongoDB** for query/read operations
- Templating engine for rendering email content
- Tracing-ready, metrics-enabled

---

## 🧱 Tech Stack

| Layer        | Tech                  |
|--------------|------------------------|
| Transport    | gRPC                   |
| CQRS Writes  | PostgreSQL             |
| CQRS Reads   | MongoDB                |
| Async Queue  | NATS                   |
| Templating   | Go `html/template`     |
| Observability| OpenTelemetry + Zap    |

---
## 🛠️ Setup

### 1. Environment Variables

Create a `.env` file or use export commands:

```env
APP_ENV=development

# Service Config
APP_SERVICE_NAME=ose-postman
APP_SERVICE_LOGGER_ENVIRONMENT=development
APP_SERVICE_LOGGER_LEVEL=info
APP_SERVICE_TRACER_ENDPOINT=localhost:4317
APP_SERVICE_TRACER_SERVICE_NAME=ose.postman
APP_SERVICE_TRACER_SAMPLE_RATIO=1.0

# Message Bus
APP_BUS_ADDRESS=nats://localhost:4222

# gRPC
APP_GRPC_PORT=20245

# PostgreSQL
APP_POSTGRES_HOST=localhost
APP_POSTGRES_PORT=5435
APP_POSTGRES_USER=fundraising
APP_POSTGRES_DATABASE=notification
APP_POSTGRES_PASSWORD=bcRqCvuAwPsbvriGXrIgSOdiuYbiGUyW
APP_POSTGRES_SSLMODE=false

# Mailer
APP_MAILER_HOST=smtp.resend.com
APP_MAILER_USERNAME=resend
APP_MAILER_PASSWORD=re_2JFNKPBB_7fuwFmWNLKvPoeALYxKaD2af
APP_MAILER_PORT=587

# MongoDB
APP_MONGO_HOST=localhost
APP_MONGO_PORT=27020
APP_MONGO_TIMEOUT=3s
APP_MONGO_USER=fundraising
APP_MONGO_PASSWORD=bcRqCvuAwPsbvriGXrIgSOdiuYbiGUyW
APP_MONGO_DATABASE=notification
```

---
# 📜 License

## MIT © Moriba Inc.

Let me know if you want this pushed into a scaffolded repo or want sample proto definitions added here as well.
