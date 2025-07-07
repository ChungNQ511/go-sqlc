Dựa trên cấu trúc thư mục của project và nội dung file `main.go`, dưới đây là phần gợi ý nội dung cho `README.md`:

---

# 🧪 Example API Service

## 🧩 Giới thiệu

Đây là project RESTful API được viết bằng Golang, sử dụng Gin làm HTTP framework, được tổ chức module hóa rõ ràng để phục vụ mục đích phát triển, testing và triển khai hệ thống backend đơn giản, mở rộng được.

## ⚙️ Công nghệ sử dụng

- **Golang**
- **Gin** – web framework
- **Swagger (swaggo)** – tài liệu API
- **PostgreSQL** – database chính
- **SQLC** – sinh code từ query SQL
- **Goose** - quản lý migration DB bằng Go/script
- **Redis** – caching
- **Docker** – đóng gói dịch vụ
- **Makefile + Shell script** – tiện ích phát triển & migration

## 🧱 Cấu trúc thư mục

```bash
.
├── db/                    # Quản lý migrations & SQLC
│   ├── migrations/        # Các file migration SQL
│   ├── queries/           # Các file query SQL
│   └── sqlc/              # File sinh bởi SQLC (db.go, models.go, querier.go, store.go)
├── docs/                  # Swagger documentation
├── modules/examples/      # Controller, router, logic ví dụ
├── pkg/                   # Thư viện hỗ trợ chia module
│   ├── clients/           # (dự phòng) external API clients
│   ├── middlewares/       # Middleware Gin (e.g., CORS)
│   └── utils/             # Các util tái sử dụng
│       ├── app/           # Kết nối DB, Redis, config
│       └── s3/            # Placeholder cho kết nối S3
├── scripts/               # Script tiện ích migration
├── main.go                # Điểm bắt đầu ứng dụng
├── sqlc.yaml              # Config cho SQLC
├── Makefile               # Lệnh build/test/tidy
└── app.env                # Biến môi trường
```

## 🚀 Khởi chạy

### 1. Cấu hình môi trường

Tạo file `.env` hoặc sửa `app.env`:

```env
SERVER_ADDRESS=8080
REDIS_HOST=localhost:6379
CORS_ORIGIN=*
SWAGGER_HOST=localhost:8080

DATABASE_DRIVER=postgres
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USERNAME=postgres
DATABASE_PASSWORD=123456789
DATABASE_NAME=example_db
```

### 2. Cài đặt dependencies

```bash
go mod tidy
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### 3. Chạy migration

```bash
make db-migrate
```

### 4. Chạy ứng dụng

```bash
go run main.go
```

### 5. Truy cập API

- Swagger: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- Health check: `GET /api/healthz`

## 📦 Build bằng Docker

```bash
docker build -t example-api .
docker run -p 8080:8080 --env-file app.env example-api
```

## 🧪 API mẫu

- `GET /api` – Trả về thông tin version
- `GET /api/healthz` – Health check
- Các route ví dụ trong module `examples`

## 📚 Swagger Documentation

Generated via [swaggo/swag](https://github.com/swaggo/swag):

```bash
swag init --parseDependency --parseInternal
```
