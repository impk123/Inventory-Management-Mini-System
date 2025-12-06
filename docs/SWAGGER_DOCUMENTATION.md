# Swagger Documentation (Full API Manual)

Inventory Management Mini System API Documentation.

## Base URL
`http://localhost:8080/api/v1`

## Authentication
ระบบใช้ **Bearer Token (JWT)** ในการยืนยันตัวตน
*   **Header**: `Authorization: Bearer <your_token>`
*   Token ได้รับจาก API `/auth/login` หรือ `/auth/register`

---

## 🔐 Authentication APIs

### 1. Register (ลงทะเบียน)
ลงทะเบียนผู้ใช้งานใหม่เข้าสู่ระบบ

*   **Method**: `POST`
*   **Endpoint**: `/auth/register`
*   **Request Body**:
    ```json
    {
      "email": "user@example.com",
      "password": "password123",
      "name": "Jack Dorson"
    }
    ```
*   **Response (201 Created)**:
    ```json
    {
      "message": "User registered successfully",
      "token": "eyJhbGciOiJIUzI1Ni...",
      "user": {
        "id": 1,
        "email": "user@example.com",
        "name": "Jack Dorson",
        "role": "user"
      }
    }
    ```

### 2. Login (เข้าสู่ระบบ)
เข้าสู่ระบบเพื่อรับ Token สำหรับใช้งาน API อื่นๆ

*   **Method**: `POST`
*   **Endpoint**: `/auth/login`
*   **Request Body**:
    ```json
    {
      "email": "user@example.com",
      "password": "password123"
    }
    ```
*   **Response (200 OK)**:
    ```json
    {
      "message": "Login successful",
      "token": "eyJhbGciOiJIUzI1Ni...",
      "user": {
        "id": 1,
        "email": "user@example.com",
        "name": "Jack Dorson",
        "role": "user"
      }
    }
    ```

### 3. Google Login
เข้าสู่ระบบผ่าน Google Account (OAuth2)

*   **Method**: `GET`
*   **Endpoint**: `/auth/google`
*   **Response**: Redirects to Google Login Page

---

## 📦 Product APIs
**Require Headers**: `Authorization: Bearer <token>`

### 1. Get All Products (ดึงรายการสินค้า)
*   **Method**: `GET`
*   **Endpoint**: `/products`
*   **Query Parameters**:
    *   `page` (int): หน้าที่ต้องการ (default: 1)
    *   `limit` (int): จำนวนต่อหน้า (default: 20)
    *   `search` (string): คำค้นหา (SKU, Name, Description)
    *   `category` (string): หมวดหมู่
*   **Response (200 OK)**:
    ```json
    {
      "products": [
        {
          "id": 1,
          "sku": "TEST-001",
          "name": "Test Product",
          "unit_price": 100,
          "quantity": 50,
          "location": "A-01"
          ...
        }
      ],
      "total": 1,
      "page": 1,
      "limit": 20
    }
    ```

### 2. Get Product By ID (ดูรายละเอียดสินค้า)
*   **Method**: `GET`
*   **Endpoint**: `/products/{id}`
*   **Path Variable**: `id` (int) - ID ของสินค้า
*   **Response (200 OK)**:
    ```json
    {
      "id": 1,
      "sku": "TEST-001",
      "name": "Test Product",
      "description": "Description...",
      "category": "General",
      "unit_price": 100,
      "cost_price": 80,
      "quantity": 50,
      "location": "A-01",
      "is_active": true
    }
    ```

### 3. Create Product (เพิ่มสินค้าใหม่)
*   **Method**: `POST`
*   **Endpoint**: `/products`
*   **Request Body**:
    ```json
    {
      "sku": "NEW-PROD-001",
      "name": "New Gaming Mouse",
      "description": "Wireless Mouse",
      "category": "Electronics",
      "unit_price": 1500,
      "cost_price": 1000,
      "quantity": 0,
      "min_quantity": 5,
      "max_quantity": 100,
      "location": "Warehouse A"
    }
    ```
*   **Response (200 OK)**:
    ```json
    {
      "message": "Product created successfully"
    }
    ```

### 4. Update Product (แก้ไขสินค้า)
*   **Method**: `PUT`
*   **Endpoint**: `/products/{id}`
*   **Request Body**: (Field ไหนไม่แก้ให้ส่งค่าเดิม หรือส่งเฉพาะ Field ที่แก้)
    ```json
    {
      "sku": "NEW-PROD-001",
      "name": "Updated Name",
      "unit_price": 1600,
       ...
    }
    ```
*   **Response (200 OK)**:
    ```json
    {
      "message": "Product updated successfully"
    }
    ```

### 5. Delete Product (ลบสินค้า)
*   **Method**: `DELETE`
*   **Endpoint**: `/products/{id}`
*   **Response (200 OK)**:
    ```json
    {
      "message": "Product deleted successfully"
    }
    ```

---

## 📊 Stock APIs
**Require Headers**: `Authorization: Bearer <token>`

### 1. Update Stock (ปรับสต็อก)
ทำรายการรับสินค้าเข้า (IN), เบิกออก (OUT), หรือปรับยอด (ADJUST)

*   **Method**: `POST`
*   **Endpoint**: `/stocks`
*   **Request Body**:
    ```json
    {
      "product_id": 1,
      "type": "IN",    // enum: "IN", "OUT", "ADJUST"
      "quantity": 10,
      "notes": "Restock from vendor"
    }
    ```
*   **Response (200 OK)**:
    ```json
    {
      "message": "Stock updated successfully"
    }
    ```

### 2. Get Current Stock (ดูสต็อกปัจจุบันทั้งหมด)
*   **Method**: `GET`
*   **Endpoint**: `/stocks`
*   **Query Parameters**:
    *   `page` (int): หน้าที่ต้องการ
    *   `limit` (int): จำนวนต่อหน้า
*   **Response (200 OK)**:
    ```json
    {
      "products": [...],
      "total": 100,
      "low_stock": [...], // รายการสินค้าที่ใกล้หมด
      "low_stock_count": 5
    }
    ```

### 3. Get Stock History By Product (ดูประวัติรายตัว)
*   **Method**: `GET`
*   **Endpoint**: `/stocks/product/{id}`
*   **Response (200 OK)**:
    ```json
    {
      "stock": [...],     // ข้อมูลปัจจุบัน
      "movements": [      // ประวัติการเคลื่อนไหวล่าสุด
        {
          "id": 1,
          "product_id": 1,
          "type": "IN",
          "quantity": 10,
          "old_quantity": 0,
          "new_quantity": 10,
          "created_at": "2023-..."
        }
      ]
    }
    ```
