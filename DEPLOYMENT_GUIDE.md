# 🚀 Deployment Guide - PSD Fiber Booking System

Complete guide for CI/CD, Docker, and structured logging implementation.

## ✅ Completed Tasks Checklist

### 1. ✅ GitHub Actions Basics
- [x] Created `.github/workflows/ci.yml` with hello-world job
- [x] Added simple echo commands to understand YAML workflow
- [x] Workflow triggers on push/PR to main branch

### 2. ✅ Run Build/Test in CI
- [x] Added `go test ./...` to CI pipeline
- [x] Tests run automatically on every commit
- [x] Binary build verification included

### 3. ✅ Dockerfile
- [x] Multi-stage Dockerfile for optimized image size
- [x] Uses Go 1.24 and Alpine Linux
- [x] Binary builds successfully

### 4. ✅ Docker Compose
- [x] Created `docker-compose.yml`
- [x] Includes API + PostgreSQL services
- [x] Health checks and networking configured
- [x] Ready for local development

### 5. ✅ CI + Docker Integration
- [x] GitHub Actions builds Docker image
- [x] Pushes to both GHCR and Docker Hub
- [x] Automated on every push to main
- [x] Tag management with metadata

### 6. ✅ Error & Logs with Logrus
- [x] Integrated logrus for structured logging
- [x] Added trace_id to all requests
- [x] JSON format for easy parsing
- [x] Error tracking with context
- [x] Ready for monitoring systems

---

## 📋 Prerequisites

Before you begin, ensure you have:

- Docker Desktop installed
- Docker Hub account (for pushing images)
- GitHub repository with Actions enabled
- Go 1.24+ (for local development)

---

## 🔧 Setup Instructions

### Step 1: Configure Docker Hub Secrets

Add these secrets to your GitHub repository (Settings → Secrets and variables → Actions):

1. **DOCKERHUB_USERNAME**: Your Docker Hub username
2. **DOCKERHUB_TOKEN**: Docker Hub access token (create at hub.docker.com/settings/security)

```bash
# To create Docker Hub token:
# 1. Go to https://hub.docker.com/settings/security
# 2. Click "New Access Token"
# 3. Name it "github-actions"
# 4. Copy the token and add to GitHub secrets
```

### Step 2: Update Docker Hub Repository Name

Edit `.github/workflows/ci.yml` line 76 to match your Docker Hub username:

```yaml
images: |
  ghcr.io/${{ github.repository }}
  YOUR_DOCKERHUB_USERNAME/psd-fiber-booking-system  # Change this
```

### Step 3: Local Development with Docker Compose

```powershell
# Clone repository
git clone https://github.com/Eursukkul/psd-fiber-booking-system.git
cd psd-fiber-booking-system

# Start services (API + Postgres)
docker compose up -d

# Check logs
docker compose logs -f api

# Stop services
docker compose down

# Stop and remove volumes
docker compose down -v
```

### Step 4: Test the API

```powershell
# Health check
curl http://localhost:3000/health

# Create a booking (with trace_id in response)
curl -X POST http://localhost:3000/api/bookings `
  -H "Content-Type: application/json" `
  -d '{"user_id":1,"service_id":101,"price":25000}'

# Get booking by ID
curl http://localhost:3000/api/bookings/1

# View Swagger docs
# Open browser: http://localhost:3000/swagger/index.html
```

---

## 📊 Structured Logging Features

### Trace ID Implementation

Every request gets a unique `trace_id` for tracking:

```json
{
  "level": "info",
  "msg": "Request completed successfully",
  "time": "2025-11-07 15:30:45",
  "trace_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "method": "POST",
  "path": "/api/bookings",
  "status": 201,
  "duration_ms": 15,
  "ip": "127.0.0.1"
}
```

### How to Use Trace ID

1. **Frontend sends trace_id** (optional):
```javascript
fetch('/api/bookings', {
  headers: {
    'X-Trace-ID': 'custom-trace-id',
    'Content-Type': 'application/json'
  }
})
```

2. **Server generates trace_id** if not provided
3. **trace_id included in all error responses**:
```json
{
  "error": "Booking not found",
  "trace_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

### Log Levels

Set via `LOG_LEVEL` environment variable:

- `debug` - Detailed debugging information
- `info` - General informational messages (default)
- `warn` - Warning messages
- `error` - Error messages only

```powershell
# Local development
$env:LOG_LEVEL="debug"
go run cmd/main.go

# Docker
docker run -e LOG_LEVEL=debug ...
```

### Example Log Output

```json
{"level":"info","msg":"Starting Fiber Booking System API","time":"2025-11-07 10:00:00"}
{"level":"info","msg":"Incoming request","time":"2025-11-07 10:00:01","trace_id":"abc123","method":"POST","path":"/api/bookings","ip":"127.0.0.1"}
{"level":"info","msg":"Creating booking","time":"2025-11-07 10:00:01","trace_id":"abc123","user_id":1,"service_id":101,"price":25000}
{"level":"info","msg":"Booking created successfully","time":"2025-11-07 10:00:01","trace_id":"abc123","booking_id":1}
{"level":"info","msg":"Request completed successfully","time":"2025-11-07 10:00:01","trace_id":"abc123","method":"POST","path":"/api/bookings","status":201,"duration_ms":45}
```

---

## 🔄 CI/CD Pipeline Flow

### Workflow Stages

1. **Hello World Job** (demonstration)
   - Echoes basic information
   - Shows workflow basics

2. **Test and Build Job**
   - Checkout code
   - Setup Go 1.24
   - Cache dependencies
   - Run `go test ./...`
   - Build binary
   - Login to GHCR & Docker Hub
   - Build & push Docker image

### Triggers

- **Push to main**: Full pipeline + image push
- **Pull Request**: Tests and build only (no push)

### Image Tags

Images are tagged with multiple versions:

- `latest` (only on main branch)
- `main` (branch name)
- `main-abc1234` (branch + commit SHA)
- Semver tags if you create releases

---

## 🐳 Docker Usage

### Build Locally

```powershell
# Build image
docker build -t psd-booking:local .

# Run container
docker run -p 3000:3000 `
  -e PORT=:3000 `
  -e JWT_SECRET=secret `
  -e LOG_LEVEL=info `
  psd-booking:local
```

### Pull from Registry

```powershell
# From GitHub Container Registry (GHCR)
docker pull ghcr.io/eursukkul/psd-fiber-booking-system:latest

# From Docker Hub (after CI push)
docker pull YOUR_USERNAME/psd-fiber-booking-system:latest

# Run pulled image
docker run -p 3000:3000 ghcr.io/eursukkul/psd-fiber-booking-system:latest
```

### Docker Compose with Custom Config

Create `docker-compose.override.yml`:

```yaml
version: '3.8'

services:
  api:
    environment:
      LOG_LEVEL: debug
      JWT_SECRET: my-custom-secret
    ports:
      - "8080:3000"  # Map to different host port
```

```powershell
# Compose uses both files automatically
docker compose up -d
```

---

## 🔍 Monitoring & Debugging

### View Logs

```powershell
# Docker Compose
docker compose logs -f api
docker compose logs -f postgres

# Docker (standalone)
docker logs -f psd-booking-api

# Filter by trace_id (if using jq)
docker logs psd-booking-api | jq 'select(.trace_id=="abc123")'
```

### Health Check

```powershell
curl http://localhost:3000/health
```

Expected response:
```json
{
  "status": "ok",
  "service": "psd-fiber-booking-system"
}
```

### Common Issues

#### Port already in use
```powershell
# Find process using port 3000
netstat -ano | findstr :3000

# Kill process by PID
taskkill /PID <PID> /F
```

#### Docker Compose fails to start
```powershell
# Check status
docker compose ps

# View detailed logs
docker compose logs

# Recreate containers
docker compose down
docker compose up --build
```

---

## 🚀 Production Deployment

### Environment Variables

Set these in production:

```bash
PORT=:3000
JWT_SECRET=<strong-random-secret>
API_KEY=<api-key>
LOG_LEVEL=info
ENV=production
DB_HOST=your-postgres-host
DB_PORT=5432
DB_USER=bookinguser
DB_PASSWORD=<strong-password>
DB_NAME=bookingdb
```

### Security Checklist

- [ ] Use strong JWT_SECRET (32+ random characters)
- [ ] Use strong database passwords
- [ ] Enable HTTPS/TLS
- [ ] Restrict CORS origins (remove `*`)
- [ ] Use secret management (AWS Secrets Manager, Azure Key Vault, etc.)
- [ ] Enable rate limiting
- [ ] Set up monitoring/alerting
- [ ] Configure log aggregation (ELK, Datadog, etc.)

### Recommended Platforms

1. **Docker Hub** → Any cloud VM
2. **AWS ECS/Fargate** → Container orchestration
3. **Azure Container Apps** → Serverless containers
4. **Google Cloud Run** → Serverless containers
5. **Kubernetes** → Self-managed or managed (EKS, AKS, GKE)

---

## 📈 Next Steps

### Monitoring Integration

The structured JSON logs are ready for:

- **ELK Stack** (Elasticsearch, Logstash, Kibana)
- **Datadog** (APM and logging)
- **New Relic** (Application monitoring)
- **Prometheus + Grafana** (Metrics and dashboards)

### Example: Shipping Logs to ELK

```yaml
# docker-compose.yml addition
services:
  api:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### Database Migration

Currently using mock repository. To integrate PostgreSQL:

1. Add database driver: `go get github.com/lib/pq`
2. Implement repository with SQL queries
3. Update docker-compose to link API to postgres
4. Add migration tool (golang-migrate, goose, etc.)

---

## 🎯 Testing the Complete Pipeline

### End-to-End Test

1. **Make a code change**
```powershell
# Edit handler/handler.go or any file
git add .
git commit -m "test: trigger CI pipeline"
git push origin main
```

2. **Watch GitHub Actions**
   - Go to repository → Actions tab
   - See workflow running
   - Check each step

3. **Verify Image Published**
```powershell
# Check GHCR
docker pull ghcr.io/eursukkul/psd-fiber-booking-system:latest

# Check Docker Hub (after adding secrets)
docker pull YOUR_USERNAME/psd-fiber-booking-system:latest
```

4. **Run and Test**
```powershell
docker run -p 3000:3000 ghcr.io/eursukkul/psd-fiber-booking-system:latest

# In another terminal
curl http://localhost:3000/health
```

---

## 📚 Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Logrus Documentation](https://github.com/sirupsen/logrus)
- [Fiber Framework](https://docs.gofiber.io/)
- [Go Testing](https://go.dev/doc/tutorial/add-a-test)

---

## 🆘 Troubleshooting

### GitHub Actions fails with "unauthorized" for Docker Hub

**Solution**: Make sure `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` secrets are set correctly.

### Trace ID not appearing in logs

**Solution**: Check that middleware is applied to all routes in `router/router.go`.

### Docker Compose postgres connection refused

**Solution**: Wait for postgres health check to pass. Use `depends_on` with `condition: service_healthy`.

---

## 📝 Summary

You now have:

✅ Automated CI/CD pipeline with GitHub Actions  
✅ Docker containerization with multi-stage builds  
✅ Docker Compose for local development (API + Postgres)  
✅ Structured logging with logrus and trace_id  
✅ Image publishing to GHCR and Docker Hub  
✅ Complete monitoring-ready setup  

**Ready for production deployment! 🎉**
