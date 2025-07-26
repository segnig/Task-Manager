# Task Manager API Documentation

**Base URL:** `/api`
**Format:** JSON
**Auth:** JWT (Token in header)

---

## 🔐 Authentication

* Protected routes **require JWT token**.
* Pass token via headers:

```
token: <JWT_TOKEN>
```

---

## 📂 User Endpoints

### 🔹 Register User

**URL:** `/api/users/register`
**Method:** `POST`
**Auth:** ❌

**Request Body:**

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "username": "johndoe123",
  "password": "securePassword",
  "user_type": "ADMIN"
}
```

**Success Response:**

```json
{
  "message": "user created successfully"
}
```

**Error Responses:**

* `400`: Malformed JSON or missing fields
* `500`: Duplicate username or user ID, validation errors

---

### 🔹 Login User

**URL:** `/api/users/login`
**Method:** `POST`
**Auth:** ❌

**Request Body:**

```json
{
  "username": "johndoe123",
  "password": "securePassword"
}
```

**Success Response:**

```json
{
  "user_id": "abc123",
  "username": "johndoe123",
  "first_name": "John",
  "last_name": "Doe",
  "token": "<JWT_TOKEN>",
  "user_type": "ADMIN"
}
```

**Error Responses:**

* `400`: User not found
* `401`: Incorrect password
* `500`: Internal error

---

### 🔹 Get All Users

**URL:** `/api/users`
**Method:** `GET`
**Auth:** ✅

**Success Response:**

```json
[
  {
    "user_id": "abc123",
    "first_name": "John",
    "last_name": "Doe",
    "username": "johndoe123"
  },
  ...
]
```

---

### 🔹 Get User by ID

**URL:** `/api/users/:user_id`
**Method:** `GET`
**Auth:** ✅

**Success Response:**

```json
{
  "user_id": "abc123",
  "first_name": "John",
  "last_name": "Doe",
  "username": "johndoe123"
}
```

---

### 🔹 Update User

**URL:** `/api/users/:user_id`
**Method:** `PUT`
**Auth:** ✅

**Request Body:**

```json
{
  "first_name": "Johnny",
  "last_name": "Doe",
  "username": "johnnydoe"
}
```

**Success Response:**

```json
{
  "message": "User Update successfully"
}
```

---

### 🔹 Delete User

**URL:** `/api/users/:user_id`
**Method:** `DELETE`
**Auth:** ✅

**Success Response:**

```json
{
  "message": "user deleted successfully"
}
```

---

## 📝 Task Endpoints

### 🔸 Create Task

**URL:** `/api/tasks`
**Method:** `POST`
**Auth:** ✅

**Request Body:**

```json
{
  "title": "Design dashboard UI",
  "description": "Use Tailwind CSS for styling"
}
```

**Success Response:**

```json
{
  "message": "Task created successfully"
}
```

**Notes:**

* `task_id`, `created_by`, `created_at`, and `updated_at` are auto-generated.

---

### 🔸 Get All Tasks

**URL:** `/api/tasks`
**Method:** `GET`
**Auth:** ❌

**Success Response:**

```json
[
  {
    "task_id": "t123",
    "title": "Design dashboard UI",
    "description": "Use Tailwind CSS",
    "created_by": "u123",
    "created_at": "2025-07-26T12:00:00Z",
    "updated_at": "2025-07-26T12:00:00Z"
  },
  ...
]
```

---

### 🔸 Get Task by ID

**URL:** `/api/tasks/:task_id`
**Method:** `GET`
**Auth:** ❌

**Success Response:**

```json
{
  "task_id": "t123",
  "title": "Design dashboard UI",
  "description": "Use Tailwind CSS",
  "created_by": "u123",
  "created_at": "2025-07-26T12:00:00Z",
  "updated_at": "2025-07-26T12:00:00Z"
}
```

---

### 🔸 Update Task

**URL:** `/api/tasks/:task_id`
**Method:** `PUT`
**Auth:** ✅

**Request Body:**

```json
{
  "title": "Update UI design",
  "description": "Change to Figma layout"
}
```

**Success Response:**

```json
{
  "message": "Task updated successfully"
}
```

**Notes:**

* Only the **creator** of the task can update it.

---

### 🔸 Delete Task

**URL:** `/api/tasks/:task_id`
**Method:** `DELETE`
**Auth:** ✅

**Success Response:**

```json
{
  "message": "Task Deleted successfully"
}
```

**Notes:**

* Only the **creator** of the task can delete it.

---

## 🧾 Models

### ✅ User

```json
{
  "user_id": "string",
  "first_name": "string",
  "last_name": "string",
  "username": "string",
  "password": "string",
  "token": "string",
  "refresh_token": "string",
  "user_type": "ADMIN | USER",
  "created_at": "ISODate",
  "updated_at": "ISODate"
}
```

---

### ✅ Task

```json
{
  "task_id": "string",
  "title": "string",
  "description": "string",
  "created_by": "user_id",
  "updated_by": "user_id",
  "created_at": "ISODate",
  "updated_at": "ISODate"
}
```

---

## 🔒 Security Rules

* The first registered user **must be an ADMIN**.
* Only task creators can **update/delete** their own tasks.
* Passwords are **hashed** before storage.
* JWT tokens are **validated** on protected routes.

---

## ✅ Status Codes

| Code | Meaning                 |
| ---- | ----------------------- |
| 200  | OK                      |
| 201  | Created                 |
| 400  | Bad Request             |
| 401  | Unauthorized            |
| 403  | Forbidden (Not Allowed) |
| 404  | Not Found               |
| 500  | Internal Server Error   |