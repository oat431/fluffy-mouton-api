# Fluffy Mouton API Documentation

> A URL shortening REST API built with Go (Fiber v3) featuring user authentication, email verification, and JWT-based access control.

**Base URL:** `/api/v1`

---

## Table of Contents

- [Standard Response Format](#standard-response-format)
- [Authentication](#authentication)
- [Endpoints](#endpoints)
  - [Health Check](#health-check)
  - [Auth - Register](#post-authregister)
  - [Auth - Login](#post-authlogin)
  - [Auth - Revoke Access](#post-authrevoke)
  - [Auth - User Details](#get-authdetail)
  - [Auth - Verify Email](#get-authverify-email)
  - [Short Link - List All](#get-short-link)
  - [Short Link - Create Random](#post-short-linkrandom)
  - [Short Link - Create Custom](#post-short-linkcustom)
  - [Redirect](#get-linktypecode)
- [Error Codes](#error-codes)
- [CORS Configuration](#cors-configuration)

---

## Standard Response Format

All endpoints return responses in the following structure:

```json
{
  "data": "<T | null>",
  "status": "SUCCESS | FAIL | ERROR",
  "error": {
    "HttpCode": 400,
    "ErrorCode": "ERROR_CODE",
    "Message": "Human-readable error message"
  }
}
```

| Field    | Type     | Description                                       |
| -------- | -------- | ------------------------------------------------- |
| `data`   | `T/null` | Response payload on success, `null` on error      |
| `status` | `string` | One of `SUCCESS`, `FAIL`, `ERROR`                 |
| `error`  | `object/null` | Error details on failure, `null` on success  |

---

## Authentication

Protected endpoints require a JWT Bearer token in the `Authorization` header:

```
Authorization: Bearer <access_token>
```

The JWT middleware validates the token and extracts the `auth_id` claim for downstream handlers.

**Unauthorized responses** (HTTP `401`) are returned when:
- The `Authorization` header is missing
- The header format is not `Bearer <token>`
- The token is invalid or expired
- The token claims cannot be parsed

---

## Endpoints

### Health Check

#### `GET /api/v1/health/check`

Returns application health status.

| Property       | Value          |
| -------------- | -------------- |
| Authentication | None           |
| Content-Type   | `application/json` |

**Response:** `200 OK` — Fiber built-in healthcheck response.

---

### `POST /auth/register`

Register a new user account. A verification email is sent upon successful registration.

| Property       | Value          |
| -------------- | -------------- |
| Authentication | None           |
| Content-Type   | `application/json` |

**Request Body:**

| Field      | Type     | Required | Validation           |
| ---------- | -------- | -------- | -------------------- |
| `username` | `string` | Yes      | —                    |
| `password` | `string` | Yes      | —                    |
| `email`    | `string` | Yes      | Must be valid email  |

```json
{
  "username": "johndoe",
  "password": "secretpassword",
  "email": "john@example.com"
}
```

**Success Response:** `201 Created`

```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john@example.com",
    "is_verified": false
  },
  "status": "SUCCESS",
  "error": null
}
```

**Error Response:** `400 Bad Request`

```json
{
  "data": null,
  "status": "ERROR",
  "error": {
    "HttpCode": 400,
    "ErrorCode": "REGISTER-01",
    "Message": "error description"
  }
}
```

---

### `POST /auth/login`

Authenticate a user and receive JWT access and refresh tokens.

| Property       | Value          |
| -------------- | -------------- |
| Authentication | None           |
| Content-Type   | `application/json` |

**Request Body:**

| Field      | Type     | Required | Validation |
| ---------- | -------- | -------- | ---------- |
| `username` | `string` | Yes      | —          |
| `password` | `string` | Yes      | —          |

```json
{
  "username": "johndoe",
  "password": "secretpassword"
}
```

**Success Response:** `200 OK`

```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  },
  "status": "SUCCESS",
  "error": null
}
```

**Error Response:** `400 Bad Request`

```json
{
  "data": null,
  "status": "ERROR",
  "error": {
    "HttpCode": 400,
    "ErrorCode": "LOGIN-01",
    "Message": "Invalid username or password"
  }
}
```

---

### `POST /auth/revoke`

Revoke a refresh token (logout).

| Property       | Value          |
| -------------- | -------------- |
| Authentication | None           |
| Content-Type   | `application/json` |

**Request Body:**

| Field           | Type     | Required | Description          |
| --------------- | -------- | -------- | -------------------- |
| `refresh_token` | `string` | Yes      | The refresh token to revoke |

```json
{
  "refresh_token": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**Success Response:** `200 OK`

```json
{
  "data": "Access revoked",
  "status": "SUCCESS",
  "error": null
}
```

**Error Responses:**

| Status | Error Code  | Message               |
| ------ | ----------- | --------------------- |
| `400`  | `REVOKE-01` | Invalid request body  |
| `500`  | `REVOKE-02` | _(dynamic error message)_ |

---

### `GET /auth/detail`

Get the authenticated user's profile details.

| Property       | Value              |
| -------------- | ------------------ |
| Authentication | **Required** (JWT) |

**Success Response:** `200 OK`

```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john@example.com",
    "is_verified": true
  },
  "status": "SUCCESS",
  "error": null
}
```

**Error Response:** `500 Internal Server Error`

```json
{
  "data": null,
  "status": "ERROR",
  "error": {
    "HttpCode": 500,
    "ErrorCode": "DETAIL-01",
    "Message": "error description"
  }
}
```

---

### `GET /auth/verify-email`

Verify a user's email address using a token sent via email.

| Property       | Value          |
| -------------- | -------------- |
| Authentication | None           |

**Query Parameters:**

| Parameter | Type     | Required | Description                  |
| --------- | -------- | -------- | ---------------------------- |
| `token`   | `string` | Yes      | Verification token from email |

**Example:** `GET /api/v1/auth/verify-email?token=abc123xyz`

**Success Response:** `200 OK`

```json
{
  "data": "Email verified successfully",
  "status": "SUCCESS",
  "error": null
}
```

**Error Responses:**

| Status | Error Code  | Message                      |
| ------ | ----------- | ---------------------------- |
| `400`  | `VERIFY-01` | Missing verification token   |
| `400`  | `VERIFY-02` | _(dynamic — invalid/expired token)_ |

---

### `GET /short-link/`

Get all short links owned by the authenticated user.

| Property       | Value              |
| -------------- | ------------------ |
| Authentication | **Required** (JWT) |

**Success Response:** `200 OK`

```json
{
  "data": [
    {
      "short_link": "abc123",
      "original_link": "https://example.com/very-long-url",
      "link_type": "RANDOM"
    },
    {
      "short_link": "my-link",
      "original_link": "https://example.com/another-url",
      "link_type": "CUSTOM"
    }
  ],
  "status": "SUCCESS",
  "error": null
}
```

**Error Response:** `500 Internal Server Error`

```json
{
  "data": null,
  "status": "ERROR",
  "error": {
    "HttpCode": 500,
    "ErrorCode": "INTERNAL_SERVER_ERROR",
    "Message": "Failed to retrieve short links"
  }
}
```

---

### `POST /short-link/random`

Create a new short link with a randomly generated code.

| Property       | Value              |
| -------------- | ------------------ |
| Authentication | **Required** (JWT) |
| Content-Type   | `application/json` |

**Request Body:**

| Field         | Type     | Required | Validation                |
| ------------- | -------- | -------- | ------------------------- |
| `url`         | `string` | Yes      | Must be a valid URL       |
| `custom_name` | `string` | No       | Alphanumeric only (ignored for random) |

```json
{
  "url": "https://example.com/my-very-long-url-that-needs-shortening"
}
```

**Success Response:** `201 Created`

```json
{
  "data": {
    "short_link": "x7kQ9m",
    "original_link": "https://example.com/my-very-long-url-that-needs-shortening",
    "link_type": "RANDOM"
  },
  "status": "SUCCESS",
  "error": null
}
```

**Error Response:** `500 Internal Server Error`

```json
{
  "data": null,
  "status": "ERROR",
  "error": {
    "HttpCode": 500,
    "ErrorCode": "INTERNAL_SERVER_ERROR",
    "Message": "Failed to create short link"
  }
}
```

---

### `POST /short-link/custom`

Create a new short link with a user-defined custom code.

| Property       | Value              |
| -------------- | ------------------ |
| Authentication | **Required** (JWT) |
| Content-Type   | `application/json` |

**Request Body:**

| Field         | Type     | Required | Validation          |
| ------------- | -------- | -------- | ------------------- |
| `url`         | `string` | Yes      | Must be a valid URL |
| `custom_name` | `string` | Yes      | Alphanumeric only   |

```json
{
  "url": "https://example.com/my-page",
  "custom_name": "mypage"
}
```

**Success Response:** `201 Created`

```json
{
  "data": {
    "short_link": "mypage",
    "original_link": "https://example.com/my-page",
    "link_type": "CUSTOM"
  },
  "status": "SUCCESS",
  "error": null
}
```

**Error Response:** `409 Conflict`

```json
{
  "data": null,
  "status": "ERROR",
  "error": {
    "HttpCode": 409,
    "ErrorCode": "SHORT_LINK_ALREADY_EXISTS",
    "Message": "Custom short link already exists"
  }
}
```

---

### `GET /:linkType/:code`

Redirect to the original URL using a short link code. This endpoint is registered at the `/api/v1` level.

| Property       | Value          |
| -------------- | -------------- |
| Authentication | None           |

**Path Parameters:**

| Parameter  | Type     | Description                                    |
| ---------- | -------- | ---------------------------------------------- |
| `linkType` | `string` | `r` for RANDOM, `c` for CUSTOM                |
| `code`     | `string` | The short link code                            |

**Example:** `GET /api/v1/r/x7kQ9m` or `GET /api/v1/c/mypage`

**Success Response:** `301 Moved Permanently` — Redirects to the original URL.

**Error Response:** `404 Not Found`

```json
{
  "data": null,
  "status": "ERROR",
  "error": {
    "HttpCode": 404,
    "ErrorCode": "SHORT_LINK_NOT_FOUND",
    "Message": "Short link not found"
  }
}
```

---

## Error Codes

| Error Code                 | HTTP Status | Endpoint        | Description                         |
| -------------------------- | ----------- | --------------- | ----------------------------------- |
| `REGISTER-01`              | 400         | POST /auth/register    | Registration failed            |
| `LOGIN-01`                 | 400         | POST /auth/login       | Invalid credentials            |
| `REVOKE-01`                | 400         | POST /auth/revoke      | Invalid request body           |
| `REVOKE-02`                | 500         | POST /auth/revoke      | Server error during revocation |
| `DETAIL-01`                | 500         | GET /auth/detail       | Failed to get user details     |
| `VERIFY-01`                | 400         | GET /auth/verify-email | Missing verification token     |
| `VERIFY-02`                | 400         | GET /auth/verify-email | Invalid or expired token       |
| `UNAUTHORIZED`             | 401         | Protected endpoints    | JWT authentication failed      |
| `INTERNAL_SERVER_ERROR`    | 500         | Short link endpoints   | Server error                   |
| `SHORT_LINK_ALREADY_EXISTS`| 409         | POST /short-link/custom| Custom code already taken      |
| `SHORT_LINK_NOT_FOUND`     | 404         | GET /:linkType/:code   | Short link does not exist      |

---

## CORS Configuration

| Setting            | Value                                            |
| ------------------ | ------------------------------------------------ |
| Allowed Origins    | `http://localhost:3000`, `http://localhost:5173`  |
| Allowed Methods    | `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`        |
| Allowed Headers    | `Origin`, `Content-Type`, `Authorization`        |
| Exposed Headers    | `Content-Length`                                  |
| Allow Credentials  | `true`                                           |
| Max Age            | 12 hours (43200 seconds)                         |
