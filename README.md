# SPD Fiber Booking System

## 📄 ภาพรวม

**SPD Fiber Booking System** is a booking system created with **Golang** and **Fiber** following the **Clean Architecture** principle. This system has features such as creating bookings, fetching booking data, canceling bookings, caching data, background tasks, generating API documentation with **Swagger**, and unit testing with **Testify**.

## 🎯 คุณสมบัติหลัก

- **Create Booking (POST /api/bookings)**

  - Create a new booking by specifying `user_id`, `service_id`, and `price`

  - Save booking data to cache immediately

  - If `price > 50,000` will check credit asynchronously and update booking status in cache

- **Get Booking by ID (GET /api/bookings/:id)**

  - Get booking data from cache or mock repository

- **Get All Bookings (GET /api/bookings)**

  - Support sorting bookings (`sort=price` or `sort=date`)

  - Support filtering bookings with high value (`high-value=true`)

- **Cancel Booking (DELETE /api/bookings/:id)**

  - Cancel booking by changing the status to "canceled"

  - Delete booking data from cache

- **Background Task**

  - Check and cancel bookings that are in the "pending" status and have been in the status for more than 5 minutes every 1 minute

- **เอกสาร API ด้วย Swagger**

  - API documentation is automatically generated using annotations in the code to specify the details of each endpoint

- **Unit Testing**

  - Cover Service Layer with **Testify** and **Mockery**

## 🛠 เทคโนโลยีที่ใช้

- **Golang** (Go 1.22)

- **Fiber**

- **Swagger** (fiber-swagger)

- **Mockery** (for testing)

- **Testify** (for Unit Testing)

- **In-Memory Cache** (can change to Redis for better performance)

- **Goroutines** (for parallel execution and background tasks)

## 📁 Project Structure

```plaintext
spd-fiber-booking-system/
├── cmd/
│   └── main.go
├── dto/
│   └── booking.go
├── handler/
│   └── booking_handler.go
├── middleware/
│   ├── auth.go
│   └── logging.go
├── router/
│   └── router.go
├── repository/
│   └── booking_repository.go
├── service/
│   └── booking_service.go
├── models/
│   └── booking.go
├── utils/
│   ├── cache.go
│   └── hash.go
├── config/
│   └── config.go
├── tests/
│   └── booking_service_test.go
├── go.mod
├── go.sum
├── .gitignore
├── .env
└── README.md
```

## 🏁 Getting Started

### 🔧 Prerequisites

- **Go** 1.22+

- **Git**

- **WSL**

### 🚀 Installation

1. **clone repository**

   ```bash
   git clone https://github.com/Eursukkul/spd-fiber-booking-system.git
   cd spd-fiber-booking-system
   ```

2. **Initialize Go Modules**

   ```bash
   go mod init github.com/Eursukkul/spd-fiber-booking-system
   ```

3. **Install Dependencies**

   ```bash
   go get github.com/gofiber/fiber/v2
   go get github.com/gofiber/swagger
   go get github.com/swaggo/swag/cmd/swag@latest
   go get github.com/stretchr/testify
   go get github.com/joho/godotenv
   go get -u github.com/swaggo/swag/cmd/swag
   go get -u github.com/swaggo/files
   go get -u github.com/swaggo/swag
   go get github.com/stretchr/testify/mock
   go get github.com/stretchr/testify/assert
   go install github.com/vektra/mockery/v2@latest
   ```

4. **Create Project Structure**

   ```bash
   mkdir -p cmd dto handler middleware router repository service models utils config tests
   touch .gitignore .env README.md
   ```

### 🔨 การเรียกใช้งานโปรเจค

1. **Create Swagger Documentation**

   ติดตั้ง Swag CLI (ถ้ายังไม่ได้ติดตั้ง):

   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. **Create Swagger Documentation**

   ```bash
   swag init
   ```

3. **Run Server**

   ```bash
   go run cmd/main.go
   ```

   The server will start at `http://localhost:3000`

4. **Access Swagger UI**

   Open the browser at [http://localhost:3000/swagger/index.html](http://localhost:3000/swagger/index.html) to view the API documentation

### 🧪 การรันการทดสอบหน่วย

ใช้คำสั่งต่อไปนี้เพื่อรันการทดสอบทั้งหมด:

```bash
go test ./...
```

## 📝 API Documentation

The system uses **Swagger** for generating API documentation. The documentation is automatically generated using annotations in the code to specify the details of each endpoint.

### 📌 ตัวอย่าง Swagger Annotations

```go
// CreateBooking สร้างการจองใหม่
// @Summary Create a new booking
// @Description Create a new booking with user_id, service_id, and price
// @Tags bookings
// @Accept  json
// @Produce  json
// @Param booking body dto.BookingRequest true "Booking"
// @Success 200 {object} dto.BookingResponse
// @Router /bookings [post]
func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
    // Implementation...
}
```

## 🛠️ Background Tasks

    ระบบมีการใช้ **Goroutines** สำหรับจัดการงานแบบ Background เช่น การตรวจสอบเครดิตและการลบการจองที่หมดอายุแล้ว

```go
func (u *bookingUsecase) BackgroundTask(wg *sync.WaitGroup) {
	ticker := time.NewTicker(1 * time.Minute)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range ticker.C {
			u.checkExpiredBookings()
		}
	}()
}

// checkExpiredBookings handles the expiration logic for pending bookings
func (u *bookingUsecase) checkExpiredBookings() {
	bookings := u.repo.GetAll()
	for _, booking := range bookings {
		if booking.Status != "pending" {
			continue
		}

		createdAt, err := time.Parse(time.RFC3339, booking.CreatedAt)
		if err != nil {
			continue
		}

		if time.Since(createdAt) > 5*time.Minute {
			u.UpdateBookingStatus(booking.ID, "canceled")
		}
	}
}
```

## 🧰 Additional Utilities

### 💾 Caching

ระบบใช้ **In-Memory Cache** สำหรับการเก็บข้อมูลการจอง เพื่อเพิ่มประสิทธิภาพในการดึงข้อมูล

```go
type InMemoryCache struct {
    store map[int]*dto.BookingResponse
    mu    sync.RWMutex
}

func NewInMemoryCache() Cache {
    return &InMemoryCache{
        store: make(map[int]*dto.BookingResponse),
    }
}

func (c *InMemoryCache) Set(id int, booking *dto.BookingResponse) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.store[id] = booking
}

func (c *InMemoryCache) Get(id int) (*dto.BookingResponse, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    booking, exists := c.store[id]
    return booking, exists
}

func (c *InMemoryCache) Delete(id int) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.store, id)
}
```

## 🧩 Clean Architecture

โปรเจคนี้ถูกจัดโครงสร้างตามหลัก **Clean Architecture** เพื่อให้มีการแยกระดับ (layers) อย่างชัดเจน ทำให้โค้ดมีความยืดหยุ่นและง่ายต่อการดูแลรักษา

### 🏛️ Layers

1. **Models:** main data structure
2. **DTOs:** Data Transfer Objects for data transfer between API
3. **Repository:** access data (use mock repository)
4. **Service:** business logic
5. **Handlers:** handle HTTP request
6. **Router:** define API routes
7. **Middleware:** handle middleware (logging, authentication)
8. **Utils:** helper tools (caching, hashing)

## 🧙‍♂️ Additional Notes

- **Mock Repository:** use for development and testing without connecting to real database
- **Goroutines:** for parallel execution and background tasks
- **Unit Testing:** cover Service Layer with **Testify**
- **Extend:** can change from In-Memory Cache to Redis for better performance in the future

---

## 📚 References

- [Fiber Documentation](https://docs.gofiber.io/)
- [Swagger Documentation](https://swagger.io/docs/)
- [Clean Architecture by Uncle Bob](https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html)

---

## 📧 Contact

หากคุณมีคำถามหรือข้อเสนอแนะเกี่ยวกับโปรเจคนี้ สามารถติดต่อได้ที่ [chalermphan.eur@gmail.com](chalermphan.eur@gmail.com)

---

## 🔐 Environment Variables

Create `.env` file in root directory:

```env
PORT=3000
JWT_SECRET=your_jwt_secret
API_KEY=your_api_key
```

## 📡 API Endpoints

| Method | Endpoint          | Description        |
| ------ | ----------------- | ------------------ |
| POST   | /api/bookings     | Create new booking |
| GET    | /api/bookings/:id | Get booking by ID  |
| GET    | /api/bookings     | Get all bookings   |
| DELETE | /api/bookings/:id | Cancel booking     |

````

## ⚠️ Error Responses

All errors follow this format:
```json
{
    "message": "Error description"
## 🚢 Deployment

### 🎯 Quick Start (5 minutes)

See **[QUICKSTART.md](QUICKSTART.md)** for immediate setup with Docker Compose!

```powershell
# Start API + PostgreSQL
docker compose up -d

# Test the API
curl http://localhost:3000/health
```

### 📖 Full Deployment Guide

For complete CI/CD, Docker Hub integration, and production deployment, see **[DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md)**.

**What's included:**
- ✅ GitHub Actions CI/CD pipeline
- ✅ Docker multi-stage build
- ✅ Docker Compose (API + PostgreSQL)
- ✅ Structured logging with trace_id (logrus + JSON)
- ✅ Auto-push to GHCR & Docker Hub
- ✅ Monitoring-ready setup

### 🚀 Deploy Options

You can deploy this app as a Docker container (recommended) or use the included GitHub Actions workflow to build and publish an image to GitHub Container Registry (GHCR).

1) Build & run locally with Docker (PowerShell example):

```powershell
# build the image (run from project root)
docker build -t fiber-booking:local .

# run the container (bind port 3000)
docker run --rm -p 3000:3000 -e PORT=:3000 -e JWT_SECRET=secret -e API_KEY=key fiber-booking:local
```

The server will be available at http://localhost:3000 and the swagger UI at /swagger/index.html

2) GitHub Actions CI (automatic build + push to GHCR)

- The workflow `.github/workflows/ci.yml` runs tests, builds the Go binary, builds a Docker image and pushes it to GitHub Container Registry (GHCR) at `ghcr.io/<ORG>/<REPO>:latest`.
- The workflow uses the repo's `GITHUB_TOKEN` with `packages: write` permission so no extra secret is required to publish to GHCR. If you prefer Docker Hub or another registry, update the workflow and add the appropriate secrets.

3) Pull & run published image from GHCR

```powershell
# example (replace ORG/REPO accordingly)
docker pull ghcr.io/<ORG>/<REPO>:latest
docker run --rm -p 3000:3000 -e PORT=:3000 ghcr.io/<ORG>/<REPO>:latest
```

Notes and tips:
- If you need a smaller image or different architecture, adjust the `Dockerfile` build stage (GOARCH/GOOS) accordingly.
- For production, prefer setting real secrets through your host or orchestration platform (Kubernetes secrets, cloud service env vars, etc.).
- To push to Docker Hub instead, replace the login/build-push steps in the CI workflow with Docker Hub credentials stored in repository secrets.

}
````
