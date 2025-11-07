# 🚀 Quick Start Guide

## ทดลองใช้งานทันที - 5 นาที

### 1️⃣ รัน Local ด้วย Docker Compose

```powershell
# เริ่มต้นทั้ง API + PostgreSQL
docker compose up -d

# ดู logs
docker compose logs -f api
```

เปิดเบราว์เซอร์:
- API: http://localhost:3000/health
- Swagger: http://localhost:3000/swagger/index.html

### 2️⃣ ทดสอบ API

```powershell
# สร้าง booking
curl -X POST http://localhost:3000/api/bookings `
  -H "Content-Type: application/json" `
  -d '{\"user_id\":1,\"service_id\":101,\"price\":25000}'

# ดู booking ทั้งหมด
curl http://localhost:3000/api/bookings

# ดู booking ตาม ID
curl http://localhost:3000/api/bookings/1
```

### 3️⃣ ดู Structured Logs

```powershell
# ดู logs แบบ JSON format พร้อม trace_id
docker compose logs api | Select-String "trace_id"
```

ตัวอย่าง log output:
```json
{
  "level": "info",
  "msg": "Creating booking",
  "trace_id": "abc-123-def",
  "user_id": 1,
  "price": 25000,
  "time": "2025-11-07 15:30:00"
}
```

### 4️⃣ ดู CI/CD Pipeline

1. Push code ไป GitHub:
```powershell
git add .
git commit -m "feat: add structured logging"
git push origin main
```

2. ไปที่ GitHub → Actions tab
3. ดู workflow รันอัตโนมัติ:
   - ✅ Run tests
   - ✅ Build binary
   - ✅ Build Docker image
   - ✅ Push to GHCR

### 5️⃣ ใช้ Image จาก Registry

```powershell
# Pull จาก GitHub Container Registry
docker pull ghcr.io/eursukkul/psd-fiber-booking-system:latest

# รัน
docker run -p 3000:3000 `
  -e PORT=:3000 `
  -e LOG_LEVEL=info `
  ghcr.io/eursukkul/psd-fiber-booking-system:latest
```

---

## 📦 สิ่งที่คุณได้รับ

### ✅ GitHub Actions Workflow
- **File**: `.github/workflows/ci.yml`
- **Features**: 
  - Hello world job (demo)
  - Auto test (`go test ./...`)
  - Auto build binary
  - Build & push Docker image

### ✅ Docker Setup
- **Dockerfile**: Multi-stage build (Go builder → Alpine runtime)
- **docker-compose.yml**: API + PostgreSQL ready to run
- **Images pushed to**: GHCR + Docker Hub

### ✅ Structured Logging
- **Library**: logrus
- **Format**: JSON
- **Features**:
  - Unique `trace_id` per request
  - Error tracking with context
  - Log levels (debug/info/warn/error)
  - Ready for monitoring tools

---

## 📋 Configuration

### Environment Variables (.env)

```env
PORT=:3000
JWT_SECRET=your_jwt_secret
API_KEY=your_api_key
LOG_LEVEL=info          # debug | info | warn | error
ENV=development         # development | production
```

### Docker Hub Setup (Optional)

ถ้าต้องการ push ไป Docker Hub:

1. เพิ่ม Secrets ใน GitHub:
   - `DOCKERHUB_USERNAME`
   - `DOCKERHUB_TOKEN`

2. แก้ไข `.github/workflows/ci.yml` บรรทัด 76:
```yaml
${{ secrets.DOCKERHUB_USERNAME }}/psd-fiber-booking-system
```

---

## 🎯 Common Commands

```powershell
# Development
go run cmd/main.go
go test ./... -v

# Docker Compose
docker compose up -d          # Start
docker compose down           # Stop
docker compose logs -f api    # View logs
docker compose ps             # Check status

# Docker Build
docker build -t psd-booking:local .
docker run -p 3000:3000 psd-booking:local

# Check logs
docker logs -f <container-id>
```

---

## 🔍 ตรวจสอบ Features

### 1. Test CI Pipeline
```powershell
# แก้ไขไฟล์อะไรก็ได้
echo "# test" >> README.md
git add .
git commit -m "test: trigger CI"
git push

# ไปดูที่ GitHub Actions
```

### 2. Test Structured Logging
```powershell
# รัน API
docker compose up -d

# สร้าง booking (จะได้ trace_id)
curl -X POST http://localhost:3000/api/bookings `
  -H "Content-Type: application/json" `
  -H "X-Trace-ID: my-custom-trace-123" `
  -d '{\"user_id\":1,\"service_id\":101,\"price\":60000}'

# ดู logs มี trace_id
docker compose logs api | Select-String "my-custom-trace-123"
```

### 3. Test Docker Image
```powershell
# Pull image
docker pull ghcr.io/eursukkul/psd-fiber-booking-system:latest

# Run
docker run -p 3000:3000 ghcr.io/eursukkul/psd-fiber-booking-system:latest

# Test
curl http://localhost:3000/health
```

---

## 🆘 Troubleshooting

### Port 3000 ถูกใช้แล้ว
```powershell
# หา process
netstat -ano | findstr :3000

# Kill process
taskkill /PID <PID> /F
```

### Docker Compose ไม่ start
```powershell
# ลบทิ้งแล้วสร้างใหม่
docker compose down -v
docker compose up --build
```

### Logs ไม่แสดง trace_id
- ตรวจสอบว่า middleware ติดตั้งใน `router/router.go`
- ตรวจสอบ LOG_LEVEL (ต้องไม่เป็น "error" เท่านั้น)

---

## 📚 เอกสารเพิ่มเติม

- **DEPLOYMENT_GUIDE.md**: คู่มือ deployment แบบละเอียด
- **README.md**: ข้อมูลโปรเจกต์และ API
- **.github/workflows/ci.yml**: CI/CD workflow
- **docker-compose.yml**: Local development setup

---

## ✨ Next Steps

1. **Add database connection**: ใช้ Postgres แทน mock repository
2. **Add monitoring**: ส่ง logs ไป ELK/Datadog
3. **Add metrics**: Prometheus + Grafana
4. **Add rate limiting**: Protect API
5. **Deploy to cloud**: AWS/Azure/GCP

---

**Happy Coding! 🎉**

ถ้ามีปัญหาหรือคำถาม สามารถดูใน `DEPLOYMENT_GUIDE.md` หรือถาม GitHub Copilot ได้เลย!
