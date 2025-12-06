# Inventory Management Mini System

ระบบ Backend API สำหรับจัดการสินค้าคงคลังขนาดเล็ก พัฒนาด้วยภาษา Go โดยเน้นโครงสร้างแบบ Clean Architecture และมีประสิทธิภาพสูงด้วยการทำ Caching

## 🚀 ฟีเจอร์หลัก (Features)

*   **Authentication**:
    *   ลงทะเบียน (Register) และเข้าสู่ระบบ (Login) ด้วย Email/Password
    *   Google OAuth2 Authentication
    *   ระบบ Token-based Authentication (JWT)
*   **Product Management**:
    *   เพิ่ม, ลบ, แก้ไข, และดูรายการสินค้า (CRUD)
    *   รองรับการค้นหาและ Pagination
*   **Stock Management**:
    *   จัดการรายการรับเข้า (Stock In) และเบิกออก (Stock Out)
    *   ดูประวัติการเคลื่อนไหวของสต็อก (Stock History)
*   **Performance**:
    *   Caching ข้อมูลสินค้าด้วย **Redis** เพื่อลดภาระของ Database
*   **Documentation**:
    *   มี API Documentation ด้วย **Swagger UI**

## 🛠️ เทคโนโลยีที่ใช้ (Tech Stack)

*   **Language**: [Go](https://go.dev/) (v1.24+)
*   **Web Framework**: [Gin Gonic](https://github.com/gin-gonic/gin)
*   **Database**: [PostgreSQL](https://www.postgresql.org/)
*   **ORM**: [GORM](https://gorm.io/)
*   **Cache**: [Redis](https://redis.io/)
*   **Documentation**: [Swagger](https://github.com/swaggo/swag)
*   **Container**: [Docker](https://www.docker.com/) & Docker Compose

## 📋 ข้อกำหนดเบื้องต้น (Prerequisites)

*   [Docker](https://www.docker.com/products/docker-desktop) และ Docker Compose (แนะนำ)
*   [Go](https://go.dev/dl/) version 1.24+ (สำหรับรันแบบ Local)

## ⚙️ การติดตั้งและรันโปรเจกต์ (Installation & Running)

### วิธีที่ 1: รันด้วย Docker (แนะนำ)

1.  **Clone โปรเจกต์**
    ```bash
    git clone https://github.com/impk123/Inventory-Management-Mini-System.git
    cd Inventory-Management-Mini-System
    ```

2.  **ตั้งค่า Environment Variables**
    ```bash
    cp .env.example .env
    ```
    *   แก้ไขไฟล์ `.env` โดยใส่ค่า Config ที่จำเป็น (ดูรายละเอียดที่หัวข้อ Environment Variables ด้านล่าง)

3.  **รันด้วย Docker Compose**
    ```bash
    docker-compose up -d --build
    ```
    *   คำสั่งนี้จะสร้าง Container สำหรับ App, PostgreSQL, และ Redis

4.  **ตรวจสอบสถานะ**
    *   Server จะรันอยู่ที่: `http://localhost:8080`
    *   Swagger Docs: `http://localhost:8080/swagger/index.html`

### วิธีที่ 2: รันแบบ Local

1.  ตรวจสอบว่ามี PostgreSQL และ Redis รันอยู่บนเครื่อง
2.  ติดตั้ง Dependencies
    ```bash
    go mod download
    ```
3.  รันโปรแกรม
    ```bash
    go run cmd/main.go
    ```

## 🔐 Environment Variables

ค่าที่ต้องกำหนดในไฟล์ `.env`:

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | Connection String สำหรับ PostgreSQL | `postgres://user:pass@localhost:5432/db` |
| `REDIS_URL` | Address ของ Redis | `localhost:6379` |
| `SERVER_PORT` | Port ที่ต้องการให้ Server รัน | `8080` |
| `JWT_SECRET` | Secret Key สำหรับสร้าง Token | `your-secret-key` |
| `GOOGLE_CLIENT_ID` | Client ID จาก Google Cloud Console | `xxxx.apps.googleusercontent.com` |
| `GOOGLE_CLIENT_SECRET` | Client Secret จาก Google Cloud Console | `xxxx` |
| `FRONTEND_URL` | URL ของ Frontend (สำหรับ Redirect หลัง Login) | `http://localhost:3000` |
| `ALLOWED_ORIGINS` | CORS configuration | `http://localhost:3000` |

## 💻 การพัฒนา (Development)

หากมีการแก้ไข API Code หรือ Comment สำหรับ Swagger สามารถ Regenerate Docs ได้ด้วยคำสั่ง:

```bash
swag init -g cmd/main.go
```
*   ต้องติดตั้ง Swag CLI ก่อน: `go install github.com/swaggo/swag/cmd/swag@latest`

## 📚 Documentation

*   **[API Documentation (Swagger Guide)](docs/API_GUIDE.md)**: วิธีการใช้งาน API และรายละเอียด Endpoint
*   **[Redis Caching Guide](docs/REDIS_GUIDE.md)**: อธิบายกลไก Caching และ Key Structure
*   **[Database Schema (ER Diagram)](docs/ER_DIAGRAM.md)**: โครงสร้างฐานข้อมูลและความสัมพันธ์

## 📂 โครงสร้างโปรเจกต์ (Project Structure)

```
Inventory-Management-Mini-System/
├── cmd/                # Entry point ของโปรแกรม main.go
├── internal/           # Business Logic หลัก
│   ├── config/         # การตั้งค่าระบบ (Env, DB)
│   ├── handlers/       # HTTP Handlers (Controller)
│   ├── models/         # Struct และ DB Models
│   ├── repositories/   # Data Access Layer
│   ├── routes/         # การกำหนด Route ของ Gin
│   └── services/       # Business Logic Layer
├── pkg/                # Packages เสริม (Utils, Middleware, Redis)
├── docs/               # Swagger Documentation Files และเอกสารอื่นๆ
├── api/                # Generated API docs
├── logs/               # Application Logs
├── docker-compose.yml  # Docker Compose Config
└── README.md           # เอกสารช่วยสอน
```

---
Developed by Impk123
