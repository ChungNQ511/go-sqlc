Dưới đây là phần `README.md` mô tả chuẩn & rõ ràng cho Makefile của bạn – bao gồm các lệnh chính, hướng dẫn dùng, và nhóm chức năng:

---

# 🛠 Makefile Commands – Project Dev Toolkit

Makefile hỗ trợ các tác vụ phát triển, migration DB, build Docker image, sinh SQL query tự động, và quản lý version release.

## 📦 Cấu trúc lệnh

---

### 🚀 Release & Docker

| Lệnh               | Mô tả                                                           |
| ------------------ | --------------------------------------------------------------- |
| `make release`     | Tăng patch version và commit file `.build-version` với git hash |
| `make build-image` | Build Docker image với version từ `.build-version`              |
| `make run-image`   | Run Docker container từ image đã build                          |
| `make push-image`  | Push Docker image lên registry                                  |

---

### 🔃 Migration DB (Goose)

| Lệnh                                       | Mô tả                                                  |
| ------------------------------------------ | ------------------------------------------------------ |
| `make create-migration MIGRATION_NAME=...` | Tạo file migration mới                                 |
| `make migrate`                             | Apply toàn bộ migration                                |
| `make db-migration-status`                 | Kiểm tra trạng thái migration                          |
| `make db-rollback`                         | Rollback nhiều phiên bản (default = 1)                 |
| `make db-rollback-version`                 | Rollback 1 bản gần nhất và sinh lại file migration tạm |
| `make db-rollback-one`                     | Rollback 1 bản gần nhất (không sinh lại file)          |
| `make remove-migration-temp`               | Xoá file migration tạm bị lỗi                          |
| `make db-migrate`                          | Chạy toàn bộ migrate, cleanup và commit lại code SQLC  |

---

### 🧠 SQLC Query Auto-Generate

| Lệnh                                                                       | Mô tả                                                 |
| -------------------------------------------------------------------------- | ----------------------------------------------------- |
| `make sqlc-db-columns TABLE_NAME=...`                                      | Lấy danh sách cột (trừ primary key)                   |
| `make sqlc-db-table-primary-key TABLE_NAME=...`                            | Lấy cột khóa chính                                    |
| `make sqlc-generate-insert-query TABLE_NAME=users [FILENAME=insert_users]` | Tạo câu lệnh `INSERT INTO ...` và lưu vào file `.sql` |
| `make sqlc-generate-update-query TABLE_NAME=users [FILENAME=update_users]` | Tạo câu lệnh `UPDATE ... WHERE pk = $1`               |
| `make sqlc-generate-delete-query TABLE_NAME=users [FILENAME=delete_users]` | Tạo câu lệnh `DELETE FROM ... WHERE pk = $1`          |

> ✅ Nếu `FILENAME` không được truyền, mặc định sẽ là `TABLE_NAME.sql`.

---

### 📦 Go module generator

| Lệnh                                                                                                     |
| -------------------------------------------------------------------------------------------------------- |
| `make go-generate-module MODULE_NAME=users [PACKAGE_DIR=users]` – sinh module Go mới theo cấu trúc chuẩn |

---

### ❓ Help

```bash
make help
```

Hiển thị danh sách lệnh có mô tả ngắn gọn.

---

## 🔧 Cấu hình

* Các biến môi trường được đọc từ `app.env`
* Yêu cầu cấu hình đúng `DATABASE_*` để dùng các lệnh liên quan tới DB / Goose / SQLC
