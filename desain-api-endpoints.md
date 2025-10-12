# Desain API Endpoints untuk Komunikasi Backend-Frontend

## 1. Pendahuluan

Desain API endpoints ini dirancang untuk mendukung komunikasi antara backend (Golang) dan frontend (Preact) dalam aplikasi WebRTC meeting. API menggunakan RESTful architecture dengan JSON sebagai format data, dan WebSocket untuk komunikasi real-time.

## 2. API Design Principles

### 2.1 RESTful Design Principles

1. **Resource-Oriented URLs**: URL harus merepresentasikan resources dengan jelas
2. **HTTP Methods**: Menggunakan HTTP methods yang sesuai (GET, POST, PUT, DELETE)
3. **Status Codes**: Menggunakan HTTP status codes yang standar
4. **JSON Format**: Menggunakan JSON untuk request dan response bodies
5. **Versioning**: API versioning untuk backward compatibility
6. **Consistent Response Format**: Format response yang konsisten di semua endpoints

### 2.2 Response Format Standard

```json
{
  "success": true,
  "data": {},
  "message": "Success message",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 2.3 Error Response Format

```json
{
  "success": false,
  "data": null,
  "message": "Error message",
  "errors": [
    {
      "field": "email",
      "message": "Email is required"
    }
  ],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 2.4 HTTP Status Codes

- **200 OK**: Request berhasil
- **201 Created**: Resource berhasil dibuat
- **204 No Content**: Request berhasil tanpa response body
- **400 Bad Request**: Request tidak valid
- **401 Unauthorized**: Authentication diperlukan
- **403 Forbidden**: Tidak memiliki permission
- **404 Not Found**: Resource tidak ditemukan
- **409 Conflict**: Conflict dengan resource yang ada
- **422 Unprocessable Entity**: Validasi gagal
- **429 Too Many Requests**: Terlalu banyak requests
- **500 Internal Server Error**: Server error
- **503 Service Unavailable**: Service tidak tersedia

## 3. Authentication Endpoints

### 3.1 Register User

**Endpoint**: `POST /api/v1/auth/register`

**Description**: Mendaftarkan pengguna baru

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "Password123!",
  "firstName": "John",
  "lastName": "Doe"
}
```

**Response** (201):
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "avatar": "",
      "status": "active",
      "emailVerified": false,
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    },
    "token": "jwt-token"
  },
  "message": "User registered successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

**Validation Errors** (422):
```json
{
  "success": false,
  "data": null,
  "message": "Validation failed",
  "errors": [
    {
      "field": "email",
      "message": "Email is required"
    },
    {
      "field": "password",
      "message": "Password must be at least 8 characters"
    }
  ],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

**Conflict Error** (409):
```json
{
  "success": false,
  "data": null,
  "message": "User with this email already exists",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 3.2 Login User

**Endpoint**: `POST /api/v1/auth/login`

**Description**: Login pengguna dengan email dan password

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "Password123!"
}
```

**Response** (200):
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "avatar": "",
      "status": "active",
      "emailVerified": false,
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    },
    "token": "jwt-token"
  },
  "message": "Login successful",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

**Authentication Error** (401):
```json
{
  "success": false,
  "data": null,
  "message": "Invalid email or password",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 3.3 Logout User

**Endpoint**: `POST /api/v1/auth/logout`

**Description**: Logout pengguna dan invalidate token

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": null,
  "message": "Logout successful",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 3.4 Refresh Token

**Endpoint**: `POST /api/v1/auth/refresh`

**Description**: Refresh JWT token

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": {
    "token": "new-jwt-token"
  },
  "message": "Token refreshed successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 3.5 Verify Email

**Endpoint**: `POST /api/v1/auth/verify-email`

**Description**: Verifikasi email pengguna

**Request Body**:
```json
{
  "token": "email-verification-token"
}
```

**Response** (200):
```json
{
  "success": true,
  "data": null,
  "message": "Email verified successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 3.6 Forgot Password

**Endpoint**: `POST /api/v1/auth/forgot-password`

**Description**: Request password reset

**Request Body**:
```json
{
  "email": "user@example.com"
}
```

**Response** (200):
```json
{
  "success": true,
  "data": null,
  "message": "Password reset email sent",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 3.7 Reset Password

**Endpoint**: `POST /api/v1/auth/reset-password`

**Description**: Reset password dengan token

**Request Body**:
```json
{
  "token": "password-reset-token",
  "password": "NewPassword123!"
}
```

**Response** (200):
```json
{
  "success": true,
  "data": null,
  "message": "Password reset successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

## 4. User Endpoints

### 4.1 Get User Profile

**Endpoint**: `GET /api/v1/users/profile`

**Description**: Mendapatkan profil pengguna yang sedang login

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "avatar": "https://example.com/avatar.jpg",
    "status": "active",
    "emailVerified": true,
    "lastLogin": "2023-01-01T00:00:00Z",
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  },
  "message": "Profile retrieved successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 4.2 Update User Profile

**Endpoint**: `PUT /api/v1/users/profile`

**Description**: Update profil pengguna

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Request Body**:
```json
{
  "firstName": "John",
  "lastName": "Doe",
  "avatar": "base64-encoded-image",
  "status": "active"
}
```

**Response** (200):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "avatar": "https://example.com/avatar.jpg",
    "status": "active",
    "emailVerified": true,
    "lastLogin": "2023-01-01T00:00:00Z",
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  },
  "message": "Profile updated successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 4.3 Change Password

**Endpoint**: `PUT /api/v1/users/password`

**Description**: Ubah password pengguna

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Request Body**:
```json
{
  "currentPassword": "CurrentPassword123!",
  "newPassword": "NewPassword123!"
}
```

**Response** (200):
```json
{
  "success": true,
  "data": null,
  "message": "Password changed successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 4.4 Get User Contacts

**Endpoint**: `GET /api/v1/users/contacts`

**Description**: Mendapatkan daftar kontak pengguna

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Query Parameters**:
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 10): Items per page
- `search` (optional): Search term

**Response** (200):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "userId": "uuid",
      "contactId": "uuid",
      "contactName": "Jane Doe",
      "contactEmail": "jane@example.com",
      "contactAvatar": "https://example.com/avatar.jpg",
      "contactStatus": "active",
      "createdAt": "2023-01-01T00:00:00Z"
    }
  ],
  "message": "Contacts retrieved successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid",
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "totalPages": 1
    }
  }
}
```

### 4.5 Add Contact

**Endpoint**: `POST /api/v1/users/contacts`

**Description**: Menambahkan kontak baru

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Request Body**:
```json
{
  "contactId": "uuid"
}
```

**Response** (201):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "userId": "uuid",
    "contactId": "uuid",
    "contactName": "Jane Doe",
    "contactEmail": "jane@example.com",
    "contactAvatar": "https://example.com/avatar.jpg",
    "contactStatus": "active",
    "createdAt": "2023-01-01T00:00:00Z"
  },
  "message": "Contact added successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 4.6 Remove Contact

**Endpoint**: `DELETE /api/v1/users/contacts/{contactId}`

**Description**: Menghapus kontak

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": null,
  "message": "Contact removed successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 4.7 Search Users

**Endpoint**: `GET /api/v1/users/search`

**Description**: Mencari pengguna untuk ditambahkan sebagai kontak

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Query Parameters**:
- `q` (required): Search query
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 10): Items per page

**Response** (200):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "email": "jane@example.com",
      "firstName": "Jane",
      "lastName": "Doe",
      "avatar": "https://example.com/avatar.jpg",
      "status": "active"
    }
  ],
  "message": "Users retrieved successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid",
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "totalPages": 1
    }
  }
}
```

### 4.8 Get User Settings

**Endpoint**: `GET /api/v1/users/settings`

**Description**: Mendapatkan pengaturan pengguna

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "userId": "uuid",
      "settingKey": "notification_enabled",
      "settingValue": "true",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    },
    {
      "id": "uuid",
      "userId": "uuid",
      "settingKey": "dark_mode",
      "settingValue": "false",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    }
  ],
  "message": "Settings retrieved successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 4.9 Update User Settings

**Endpoint**: `PUT /api/v1/users/settings`

**Description**: Update pengaturan pengguna

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Request Body**:
```json
{
  "settings": [
    {
      "settingKey": "notification_enabled",
      "settingValue": "false"
    },
    {
      "settingKey": "dark_mode",
      "settingValue": "true"
    }
  ]
}
```

**Response** (200):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "userId": "uuid",
      "settingKey": "notification_enabled",
      "settingValue": "false",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    },
    {
      "id": "uuid",
      "userId": "uuid",
      "settingKey": "dark_mode",
      "settingValue": "true",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    }
  ],
  "message": "Settings updated successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

## 5. Room Endpoints

### 5.1 Create Room

**Endpoint**: `POST /api/v1/rooms`

**Description**: Membuat ruang meeting baru

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Request Body**:
```json
{
  "name": "Team Meeting",
  "description": "Weekly team sync",
  "password": "password123",
  "maxUsers": 10,
  "scheduledAt": "2023-01-01T10:00:00Z"
}
```

**Response** (201):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Team Meeting",
    "description": "Weekly team sync",
    "hostId": "uuid",
    "hostName": "John Doe",
    "password": "password123",
    "maxUsers": 10,
    "status": "active",
    "scheduledAt": "2023-01-01T10:00:00Z",
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  },
  "message": "Room created successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 5.2 Get Rooms

**Endpoint**: `GET /api/v1/rooms`

**Description**: Mendapatkan daftar ruang meeting

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Query Parameters**:
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 10): Items per page
- `status` (optional): Filter by status (active, ended, cancelled)
- `hostId` (optional): Filter by host ID

**Response** (200):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "Team Meeting",
      "description": "Weekly team sync",
      "hostId": "uuid",
      "hostName": "John Doe",
      "maxUsers": 10,
      "status": "active",
      "scheduledAt": "2023-01-01T10:00:00Z",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z",
      "currentParticipants": 3
    }
  ],
  "message": "Rooms retrieved successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid",
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "totalPages": 1
    }
  }
}
```

### 5.3 Get Room

**Endpoint**: `GET /api/v1/rooms/{roomId}`

**Description**: Mendapatkan detail ruang meeting

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Team Meeting",
    "description": "Weekly team sync",
    "hostId": "uuid",
    "hostName": "John Doe",
    "password": "password123",
    "maxUsers": 10,
    "status": "active",
    "scheduledAt": "2023-01-01T10:00:00Z",
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z",
    "currentParticipants": 3,
    "participants": [
      {
        "userId": "uuid",
        "userName": "John Doe",
        "userAvatar": "https://example.com/avatar.jpg",
        "role": "host",
        "joinedAt": "2023-01-01T00:00:00Z"
      },
      {
        "userId": "uuid",
        "userName": "Jane Doe",
        "userAvatar": "https://example.com/avatar.jpg",
        "role": "participant",
        "joinedAt": "2023-01-01T00:05:00Z"
      }
    ]
  },
  "message": "Room retrieved successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 5.4 Update Room

**Endpoint**: `PUT /api/v1/rooms/{roomId}`

**Description**: Update ruang meeting

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Request Body**:
```json
{
  "name": "Updated Team Meeting",
  "description": "Updated weekly team sync",
  "password": "newpassword123",
  "maxUsers": 15,
  "scheduledAt": "2023-01-01T11:00:00Z"
}
```

**Response** (200):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Updated Team Meeting",
    "description": "Updated weekly team sync",
    "hostId": "uuid",
    "hostName": "John Doe",
    "password": "newpassword123",
    "maxUsers": 15,
    "status": "active",
    "scheduledAt": "2023-01-01T11:00:00Z",
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  },
  "message": "Room updated successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 5.5 Delete Room

**Endpoint**: `DELETE /api/v1/rooms/{roomId}`

**Description**: Menghapus ruang meeting

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": null,
  "message": "Room deleted successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 5.6 Join Room

**Endpoint**: `POST /api/v1/rooms/{roomId}/join`

**Description**: Bergabung ke ruang meeting

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Request Body**:
```json
{
  "password": "password123"
}
```

**Response** (200):
```json
{
  "success": true,
  "data": {
    "roomId": "uuid",
    "userId": "uuid",
    "role": "participant",
    "joinedAt": "2023-01-01T00:00:00Z"
  },
  "message": "Joined room successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 5.7 Leave Room

**Endpoint**: `POST /api/v1/rooms/{roomId}/leave`

**Description**: Meninggalkan ruang meeting

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": null,
  "message": "Left room successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 5.8 End Room

**Endpoint**: `POST /api/v1/rooms/{roomId}/end`

**Description**: Mengakhiri ruang meeting

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Team Meeting",
    "status": "ended",
    "endedAt": "2023-01-01T00:00:00Z"
  },
  "message": "Room ended successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 5.9 Get Room Participants

**Endpoint**: `GET /api/v1/rooms/{roomId}/participants`

**Description**: Mendapatkan daftar peserta dalam ruang meeting

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": [
    {
      "userId": "uuid",
      "userName": "John Doe",
      "userAvatar": "https://example.com/avatar.jpg",
      "role": "host",
      "joinedAt": "2023-01-01T00:00:00Z",
      "isInRoom": true
    },
    {
      "userId": "uuid",
      "userName": "Jane Doe",
      "userAvatar": "https://example.com/avatar.jpg",
      "role": "participant",
      "joinedAt": "2023-01-01T00:05:00Z",
      "isInRoom": true
    }
  ],
  "message": "Participants retrieved successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 5.10 Get Room History

**Endpoint**: `GET /api/v1/rooms/{roomId}/history`

**Description**: Mendapatkan riwayat ruang meeting

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Query Parameters**:
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 10): Items per page

**Response** (200):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "roomId": "uuid",
      "userId": "uuid",
      "userName": "John Doe",
      "joinedAt": "2023-01-01T00:00:00Z",
      "leftAt": "2023-01-01T01:00:00Z",
      "duration": "1h 0m 0s",
      "recordingUrl": "https://example.com/recording.mp4"
    }
  ],
  "message": "Room history retrieved successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid",
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "totalPages": 1
    }
  }
}
```

### 5.11 Get Room Messages

**Endpoint**: `GET /api/v1/rooms/{roomId}/messages`

**Description**: Mendapatkan pesan dalam ruang meeting

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Query Parameters**:
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 50): Items per page

**Response** (200):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "roomId": "uuid",
      "userId": "uuid",
      "userName": "John Doe",
      "userAvatar": "https://example.com/avatar.jpg",
      "message": "Hello everyone!",
      "messageType": "text",
      "createdAt": "2023-01-01T00:00:00Z"
    },
    {
      "id": "uuid",
      "roomId": "uuid",
      "userId": "uuid",
      "userName": "Jane Doe",
      "userAvatar": "https://example.com/avatar.jpg",
      "message": "Hi John!",
      "messageType": "text",
      "createdAt": "2023-01-01T00:01:00Z"
    }
  ],
  "message": "Messages retrieved successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid",
    "pagination": {
      "page": 1,
      "limit": 50,
      "total": 2,
      "totalPages": 1
    }
  }
}
```

### 5.12 Get Room Settings

**Endpoint**: `GET /api/v1/rooms/{roomId}/settings`

**Description**: Mendapatkan pengaturan ruang meeting

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "roomId": "uuid",
      "settingKey": "recording_enabled",
      "settingValue": "true",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    },
    {
      "id": "uuid",
      "roomId": "uuid",
      "settingKey": "waiting_room",
      "settingValue": "false",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    }
  ],
  "message": "Room settings retrieved successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 5.13 Update Room Settings

**Endpoint**: `PUT /api/v1/rooms/{roomId}/settings`

**Description**: Update pengaturan ruang meeting

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Request Body**:
```json
{
  "settings": [
    {
      "settingKey": "recording_enabled",
      "settingValue": "false"
    },
    {
      "settingKey": "waiting_room",
      "settingValue": "true"
    }
  ]
}
```

**Response** (200):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "roomId": "uuid",
      "settingKey": "recording_enabled",
      "settingValue": "false",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    },
    {
      "id": "uuid",
      "roomId": "uuid",
      "settingKey": "waiting_room",
      "settingValue": "true",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    }
  ],
  "message": "Room settings updated successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

## 6. WebSocket Endpoints

### 6.1 WebSocket Connection

**Endpoint**: `ws://localhost:8080/ws?roomId={roomId}&userId={userId}&token={jwt-token}`

**Description**: Koneksi WebSocket untuk komunikasi real-time

### 6.2 WebSocket Message Types

#### 6.2.1 Join Room

**Client to Server**:
```json
{
  "type": "join_room",
  "roomId": "uuid",
  "userId": "uuid",
  "payload": {
    "username": "John Doe"
  }
}
```

**Server to Client (Broadcast)**:
```json
{
  "type": "user_joined",
  "roomId": "uuid",
  "userId": "uuid",
  "username": "John Doe",
  "payload": {
    "message": "User joined the room"
  },
  "timestamp": "2023-01-01T00:00:00Z"
}
```

#### 6.2.2 Leave Room

**Client to Server**:
```json
{
  "type": "leave_room",
  "roomId": "uuid",
  "userId": "uuid",
  "payload": {
    "username": "John Doe"
  }
}
```

**Server to Client (Broadcast)**:
```json
{
  "type": "user_left",
  "roomId": "uuid",
  "userId": "uuid",
  "username": "John Doe",
  "payload": {
    "message": "User left the room"
  },
  "timestamp": "2023-01-01T00:00:00Z"
}
```

#### 6.2.3 Chat Message

**Client to Server**:
```json
{
  "type": "chat",
  "roomId": "uuid",
  "userId": "uuid",
  "payload": {
    "message": "Hello everyone!"
  }
}
```

**Server to Client (Broadcast)**:
```json
{
  "type": "chat",
  "roomId": "uuid",
  "userId": "uuid",
  "username": "John Doe",
  "payload": {
    "message": "Hello everyone!"
  },
  "timestamp": "2023-01-01T00:00:00Z"
}
```

#### 6.2.4 WebRTC Offer

**Client to Server**:
```json
{
  "type": "offer",
  "roomId": "uuid",
  "userId": "uuid",
  "targetUserId": "uuid",
  "payload": {
    "sdp": "webrtc-sdp-offer"
  }
}
```

**Server to Client (Send to target user)**:
```json
{
  "type": "offer",
  "roomId": "uuid",
  "userId": "uuid",
  "username": "John Doe",
  "payload": {
    "sdp": "webrtc-sdp-offer"
  },
  "timestamp": "2023-01-01T00:00:00Z"
}
```

#### 6.2.5 WebRTC Answer

**Client to Server**:
```json
{
  "type": "answer",
  "roomId": "uuid",
  "userId": "uuid",
  "targetUserId": "uuid",
  "payload": {
    "sdp": "webrtc-sdp-answer"
  }
}
```

**Server to Client (Send to target user)**:
```json
{
  "type": "answer",
  "roomId": "uuid",
  "userId": "uuid",
  "username": "John Doe",
  "payload": {
    "sdp": "webrtc-sdp-answer"
  },
  "timestamp": "2023-01-01T00:00:00Z"
}
```

#### 6.2.6 WebRTC ICE Candidate

**Client to Server**:
```json
{
  "type": "ice-candidate",
  "roomId": "uuid",
  "userId": "uuid",
  "targetUserId": "uuid",
  "payload": {
    "candidate": {
      "candidate": "candidate-string",
      "sdpMid": "sdp-mid",
      "sdpMLineIndex": 0
    }
  }
}
```

**Server to Client (Send to target user)**:
```json
{
  "type": "ice-candidate",
  "roomId": "uuid",
  "userId": "uuid",
  "username": "John Doe",
  "payload": {
    "candidate": {
      "candidate": "candidate-string",
      "sdpMid": "sdp-mid",
      "sdpMLineIndex": 0
    }
  },
  "timestamp": "2023-01-01T00:00:00Z"
}
```

#### 6.2.7 Mute/Unmute

**Client to Server**:
```json
{
  "type": "mute",
  "roomId": "uuid",
  "userId": "uuid",
  "payload": {
    "muted": true
  }
}
```

**Server to Client (Broadcast)**:
```json
{
  "type": "mute",
  "roomId": "uuid",
  "userId": "uuid",
  "username": "John Doe",
  "payload": {
    "muted": true
  },
  "timestamp": "2023-01-01T00:00:00Z"
}
```

#### 6.2.8 Video On/Off

**Client to Server**:
```json
{
  "type": "video",
  "roomId": "uuid",
  "userId": "uuid",
  "payload": {
    "videoEnabled": false
  }
}
```

**Server to Client (Broadcast)**:
```json
{
  "type": "video",
  "roomId": "uuid",
  "userId": "uuid",
  "username": "John Doe",
  "payload": {
    "videoEnabled": false
  },
  "timestamp": "2023-01-01T00:00:00Z"
}
```

#### 6.2.9 Screen Share

**Client to Server**:
```json
{
  "type": "screen_share",
  "roomId": "uuid",
  "userId": "uuid",
  "payload": {
    "enabled": true
  }
}
```

**Server to Client (Broadcast)**:
```json
{
  "type": "screen_share",
  "roomId": "uuid",
  "userId": "uuid",
  "username": "John Doe",
  "payload": {
    "enabled": true
  },
  "timestamp": "2023-01-01T00:00:00Z"
}
```

#### 6.2.10 Recording

**Client to Server**:
```json
{
  "type": "recording",
  "roomId": "uuid",
  "userId": "uuid",
  "payload": {
    "enabled": true
  }
}
```

**Server to Client (Broadcast)**:
```json
{
  "type": "recording",
  "roomId": "uuid",
  "userId": "uuid",
  "username": "John Doe",
  "payload": {
    "enabled": true,
    "recordingUrl": "https://example.com/recording.mp4"
  },
  "timestamp": "2023-01-01T00:00:00Z"
}
```

#### 6.2.11 Error

**Server to Client**:
```json
{
  "type": "error",
  "roomId": "uuid",
  "userId": "uuid",
  "payload": {
    "message": "Error message",
    "code": "ERROR_CODE"
  },
  "timestamp": "2023-01-01T00:00:00Z"
}
```

## 7. Notification Endpoints

### 7.1 Get Notifications

**Endpoint**: `GET /api/v1/notifications`

**Description**: Mendapatkan notifikasi pengguna

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Query Parameters**:
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 10): Items per page
- `isRead` (optional): Filter by read status

**Response** (200):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "type": "room_invitation",
      "title": "Room Invitation",
      "message": "John Doe invited you to join 'Team Meeting'",
      "isRead": false,
      "createdAt": "2023-01-01T00:00:00Z"
    },
    {
      "id": "uuid",
      "type": "meeting_reminder",
      "title": "Meeting Reminder",
      "message": "Your meeting 'Team Meeting' will start in 15 minutes",
      "isRead": true,
      "createdAt": "2023-01-01T00:00:00Z"
    }
  ],
  "message": "Notifications retrieved successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid",
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 2,
      "totalPages": 1
    }
  }
}
```

### 7.2 Mark Notification as Read

**Endpoint**: `PUT /api/v1/notifications/{notificationId}/read`

**Description**: Menandai notifikasi sebagai telah dibaca

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": null,
  "message": "Notification marked as read",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 7.3 Mark All Notifications as Read

**Endpoint**: `PUT /api/v1/notifications/read-all`

**Description**: Menandai semua notifikasi sebagai telah dibaca

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": null,
  "message": "All notifications marked as read",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 7.4 Delete Notification

**Endpoint**: `DELETE /api/v1/notifications/{notificationId}`

**Description**: Menghapus notifikasi

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": null,
  "message": "Notification deleted successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

## 8. WebRTC Endpoints

### 8.1 Generate ICE Servers

**Endpoint**: `GET /api/v1/webrtc/ice-servers`

**Description**: Mendapatkan konfigurasi ICE servers

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": {
    "iceServers": [
      {
        "urls": [
          "stun:stun.l.google.com:19302",
          "stun:stun1.l.google.com:19302"
        ]
      },
      {
        "urls": "turn:your-turn-server.com:3478",
        "username": "username",
        "credential": "credential"
      }
    ]
  },
  "message": "ICE servers retrieved successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 8.2 Start Recording

**Endpoint**: `POST /api/v1/webrtc/rooms/{roomId}/recording/start`

**Description**: Memulai recording meeting

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": {
    "recordingId": "uuid",
    "recordingUrl": "https://example.com/recording.mp4",
    "startedAt": "2023-01-01T00:00:00Z"
  },
  "message": "Recording started successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 8.3 Stop Recording

**Endpoint**: `POST /api/v1/webrtc/rooms/{roomId}/recording/stop`

**Description**: Menghentikan recording meeting

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": {
    "recordingId": "uuid",
    "recordingUrl": "https://example.com/recording.mp4",
    "duration": "1h 0m 0s",
    "stoppedAt": "2023-01-01T01:00:00Z"
  },
  "message": "Recording stopped successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 8.4 Get Recording

**Endpoint**: `GET /api/v1/webrtc/recordings/{recordingId}`

**Description**: Mendapatkan informasi recording

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "roomId": "uuid",
    "roomName": "Team Meeting",
    "recordingUrl": "https://example.com/recording.mp4",
    "duration": "1h 0m 0s",
    "startedAt": "2023-01-01T00:00:00Z",
    "stoppedAt": "2023-01-01T01:00:00Z",
    "createdAt": "2023-01-01T00:00:00Z"
  },
  "message": "Recording retrieved successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 8.5 Get Recordings

**Endpoint**: `GET /api/v1/webrtc/recordings`

**Description**: Mendapatkan daftar recording

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Query Parameters**:
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 10): Items per page
- `roomId` (optional): Filter by room ID

**Response** (200):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "roomId": "uuid",
      "roomName": "Team Meeting",
      "recordingUrl": "https://example.com/recording.mp4",
      "duration": "1h 0m 0s",
      "startedAt": "2023-01-01T00:00:00Z",
      "stoppedAt": "2023-01-01T01:00:00Z",
      "createdAt": "2023-01-01T00:00:00Z"
    }
  ],
  "message": "Recordings retrieved successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid",
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "totalPages": 1
    }
  }
}
```

## 9. Admin Endpoints

### 9.1 Get Users (Admin)

**Endpoint**: `GET /api/v1/admin/users`

**Description**: Mendapatkan daftar pengguna (admin only)

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Query Parameters**:
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 10): Items per page
- `status` (optional): Filter by status
- `search` (optional): Search term

**Response** (200):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "email": "user@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "avatar": "https://example.com/avatar.jpg",
      "status": "active",
      "emailVerified": true,
      "lastLogin": "2023-01-01T00:00:00Z",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    }
  ],
  "message": "Users retrieved successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid",
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "totalPages": 1
    }
  }
}
```

### 9.2 Update User Status (Admin)

**Endpoint**: `PUT /api/v1/admin/users/{userId}/status`

**Description**: Update status pengguna (admin only)

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Request Body**:
```json
{
  "status": "banned"
}
```

**Response** (200):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "banned",
    "updatedAt": "2023-01-01T00:00:00Z"
  },
  "message": "User status updated successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

### 9.3 Get Rooms (Admin)

**Endpoint**: `GET /api/v1/admin/rooms`

**Description**: Mendapatkan daftar ruang meeting (admin only)

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Query Parameters**:
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 10): Items per page
- `status` (optional): Filter by status
- `hostId` (optional): Filter by host ID

**Response** (200):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "Team Meeting",
      "description": "Weekly team sync",
      "hostId": "uuid",
      "hostName": "John Doe",
      "maxUsers": 10,
      "status": "active",
      "scheduledAt": "2023-01-01T10:00:00Z",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z",
      "currentParticipants": 3
    }
  ],
  "message": "Rooms retrieved successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid",
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "totalPages": 1
    }
  }
}
```

### 9.4 Get System Statistics (Admin)

**Endpoint**: `GET /api/v1/admin/statistics`

**Description**: Mendapatkan statistik sistem (admin only)

**Request Headers**:
```
Authorization: Bearer jwt-token
```

**Response** (200):
```json
{
  "success": true,
  "data": {
    "users": {
      "total": 1000,
      "active": 850,
      "inactive": 100,
      "banned": 50
    },
    "rooms": {
      "total": 500,
      "active": 100,
      "ended": 400
    },
    "meetings": {
      "total": 2000,
      "today": 50,
      "thisWeek": 300,
      "thisMonth": 1000
    },
    "recordings": {
      "total": 500,
      "totalSize": "50GB",
      "today": 10,
      "thisWeek": 50,
      "thisMonth": 200
    }
  },
  "message": "Statistics retrieved successfully",
  "errors": [],
  "meta": {
    "timestamp": "2023-01-01T00:00:00Z",
    "requestId": "uuid"
  }
}
```

## 10. API Security

### 10.1 Authentication

- **JWT (JSON Web Token)** untuk autentikasi stateless
- Token expiration: 24 jam untuk access token, 7 hari untuk refresh token
- Token disimpan di HTTP Authorization header: `Authorization: Bearer <token>`
- Refresh token mechanism untuk memperpanjang sesi tanpa login ulang

### 10.2 Authorization

- **Role-based access control (RBAC)** untuk mengelola permission
- Row-level security di database untuk membatasi akses data
- Middleware untuk validasi permission di setiap endpoint

### 10.3 Rate Limiting

- Rate limiting per endpoint dan per user
- Redis untuk menyimpan counter rate limiting
- Response headers dengan rate limit info:
  ```
  X-RateLimit-Limit: 100
  X-RateLimit-Remaining: 99
  X-RateLimit-Reset: 1640995200
  ```

### 10.4 Input Validation

- Validasi input di semua endpoint
- Sanitasi input untuk mencegah injection attacks
- Custom error messages untuk validasi errors

### 10.5 CORS

- CORS configuration untuk mengizinkan akses dari frontend domain
- Preflight requests handling untuk complex requests

### 10.6 HTTPS

- HTTPS untuk semua API endpoints di production
- HSTS (HTTP Strict Transport Security) header
- Certificate pinning untuk enhanced security

## 11. API Documentation

### 11.1 OpenAPI/Swagger Documentation

API documentation akan di-generate menggunakan OpenAPI/Swagger specification:

```yaml
openapi: 3.0.0
info:
  title: WebRTC Meeting API
  description: API for WebRTC meeting application
  version: 1.0.0
  contact:
    name: API Support
    email: support@example.com
servers:
  - url: https://api.example.com/api/v1
    description: Production server
  - url: https://staging-api.example.com/api/v1
    description: Staging server
  - url: http://localhost:8080/api/v1
    description: Development server
paths:
  /auth/login:
    post:
      tags:
        - Authentication
      summary: Login user
      description: Authenticate user with email and password
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - email
                - password
              properties:
                email:
                  type: string
                  format: email
                  example: user@example.com
                password:
                  type: string
                  format: password
                  example: Password123!
      responses:
        '200':
          description: Login successful
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: object
                    properties:
                      user:
                        $ref: '#/components/schemas/User'
                      token:
                        type: string
                        example: jwt-token
                  message:
                    type: string
                    example: Login successful
                  errors:
                    type: array
                    items:
                      type: object
                    example: []
                  meta:
                    $ref: '#/components/schemas/Meta'
        '401':
          $ref: '#/components/responses/UnauthorizedError'
        '422':
          $ref: '#/components/responses/ValidationError'
components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: string
          format: uuid
          example: 550e8400-e29b-41d4-a716-446655440000
        email:
          type: string
          format: email
          example: user@example.com
        firstName:
          type: string
          example: John
        lastName:
          type: string
          example: Doe
        avatar:
          type: string
          format: uri
          example: https://example.com/avatar.jpg
        status:
          type: string
          enum: [active, inactive, banned]
          example: active
        emailVerified:
          type: boolean
          example: true
        lastLogin:
          type: string
          format: date-time
          example: 2023-01-01T00:00:00Z
        createdAt:
          type: string
          format: date-time
          example: 2023-01-01T00:00:00Z
        updatedAt:
          type: string
          format: date-time
          example: 2023-01-01T00:00:00Z
    Meta:
      type: object
      properties:
        timestamp:
          type: string
          format: date-time
          example: 2023-01-01T00:00:00Z
        requestId:
          type: string
          format: uuid
          example: 550e8400-e29b-41d4-a716-446655440000
        pagination:
          type: object
          properties:
            page:
              type: integer
              example: 1
            limit:
              type: integer
              example: 10
            total:
              type: integer
              example: 100
            totalPages:
              type: integer
              example: 10
  responses:
    UnauthorizedError:
      description: Authentication failed
      content:
        application/json:
          schema:
            type: object
            properties:
              success:
                type: boolean
                example: false
              data:
                type: object
                nullable: true
                example: null
              message:
                type: string
                example: Invalid email or password
              errors:
                type: array
                items:
                  type: object
                example: []
              meta:
                $ref: '#/components/schemas/Meta'
    ValidationError:
      description: Validation failed
      content:
        application/json:
          schema:
            type: object
            properties:
              success:
                type: boolean
                example: false
              data:
                type: object
                nullable: true
                example: null
              message:
                type: string
                example: Validation failed
              errors:
                type: array
                items:
                  type: object
                  properties:
                    field:
                      type: string
                      example: email
                    message:
                      type: string
                      example: Email is required
                example:
                  - field: email
                    message: Email is required
                  - field: password
                    message: Password must be at least 8 characters
              meta:
                $ref: '#/components/schemas/Meta'
```

### 11.2 Interactive API Documentation

Menggunakan **Swagger UI** atau **Redoc** untuk interactive API documentation:

```bash
# Install swagger-ui-express
npm install swagger-ui-express

# Serve API documentation
app.use('/api-docs', swaggerUi.serve, swaggerUi.setup(swaggerSpec));
```

### 11.3 API Client Generation

Menggunakan **OpenAPI Generator** untuk mengenerate API client untuk frontend:

```bash
# Generate TypeScript client
openapi-generator generate \
  -i ./openapi.yaml \
  -g typescript-fetch \
  -o ./src/services/api/generated
```

## 12. API Testing

### 12.1 Unit Testing

Menggunakan **Go testing** dan **testify** untuk unit testing API endpoints:

```go
// main_test.go
package main

import (
  "bytes"
  "encoding/json"
  "net/http"
  "net/http/httptest"
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
  // Setup
  router := setupRouter()

  // Test data
  loginData := map[string]interface{}{
    "email":    "test@example.com",
    "password": "password123",
  }
  jsonData, _ := json.Marshal(loginData)

  // Create request
  req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
  req.Header.Set("Content-Type", "application/json")

  // Record response
  w := httptest.NewRecorder()
  router.ServeHTTP(w, req)

  // Assertions
  assert.Equal(t, http.StatusOK, w.Code)
  
  var response map[string]interface{}
  json.Unmarshal(w.Body.Bytes(), &response)
  
  assert.Equal(t, true, response["success"])
  assert.NotNil(t, response["data"])
  assert.NotNil(t, response["data"].(map[string]interface{})["token"])
}
```

### 12.2 Integration Testing

Menggunakan **testcontainers** untuk integration testing dengan database:

```go
// integration_test.go
package main

import (
  "context"
  "database/sql"
  "testing"
  "time"

  "github.com/stretchr/testify/assert"
  "github.com/testcontainers/testcontainers-go"
  "github.com/testcontainers/testcontainers-go/wait"
)

func TestUserEndpointsWithDatabase(t *testing.T) {
  // Setup PostgreSQL container
  ctx := context.Background()
  req := testcontainers.ContainerRequest{
    Image:        "postgres:13",
    ExposedPorts: []string{"5432/tcp"},
    Env: map[string]string{
      "POSTGRES_PASSWORD": "password",
      "POSTGRES_USER":     "user",
      "POSTGRES_DB":       "testdb",
    },
    WaitingFor: wait.ForLog("database system is ready to accept connections").
      WithOccurrence(2).
      WithStartupTimeout(5 * time.Second),
  }
  
  postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
    ContainerRequest: req,
    Started:          true,
  })
  assert.NoError(t, err)
  
  defer postgresContainer.Terminate(ctx)
  
  // Get database connection
  host, _ := postgresContainer.Host(ctx)
  port, _ := postgresContainer.MappedPort(ctx, "5432")
  
  db, err := sql.Open("postgres", fmt.Sprintf("postgres://user:password@%s:%s/testdb?sslmode=disable", host, port.Port()))
  assert.NoError(t, err)
  defer db.Close()
  
  // Run migrations
  err = runMigrations(db)
  assert.NoError(t, err)
  
  // Setup router with database connection
  router := setupRouterWithDB(db)
  
  // Test user registration
  registerData := map[string]interface{}{
    "email":     "test@example.com",
    "password":  "password123",
    "firstName": "Test",
    "lastName":  "User",
  }
  jsonData, _ := json.Marshal(registerData)
  
  req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
  req.Header.Set("Content-Type", "application/json")
  
  w := httptest.NewRecorder()
  router.ServeHTTP(w, req)
  
  assert.Equal(t, http.StatusCreated, w.Code)
  
  var response map[string]interface{}
  json.Unmarshal(w.Body.Bytes(), &response)
  
  assert.Equal(t, true, response["success"])
  assert.NotNil(t, response["data"])
  
  // Verify user in database
  var count int
  err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", "test@example.com").Scan(&count)
  assert.NoError(t, err)
  assert.Equal(t, 1, count)
}
```

### 12.3 API Contract Testing

Menggunakan **Pact** untuk API contract testing:

```javascript
// consumer.test.js
const { Pact } = require('@pact-foundation/pact');
const { like, eachLike } = require('@pact-foundation/pact/dsl/matchers');

describe('API Contract Test', () => {
  const provider = new Pact({
    consumer: 'WebRTC Meeting Frontend',
    provider: 'WebRTC Meeting API',
    port: 8080,
    log: './logs/pact.log',
    dir: './pacts',
  });

  beforeAll(() => provider.setup());
  afterEach(() => provider.verify());
  afterAll(() => provider.finalize());

  describe('Login API', () => {
    test('should return user data and token on successful login', async () => {
      await provider.addInteraction({
        state: 'User exists',
        uponReceiving: 'a request to login',
        withRequest: {
          method: 'POST',
          path: '/api/v1/auth/login',
          headers: {
            'Content-Type': 'application/json',
          },
          body: {
            email: 'test@example.com',
            password: 'password123',
          },
        },
        willRespondWith: {
          status: 200,
          headers: {
            'Content-Type': 'application/json',
          },
          body: {
            success: true,
            data: {
              user: {
                id: like('550e8400-e29b-41d4-a716-446655440000'),
                email: 'test@example.com',
                firstName: 'Test',
                lastName: 'User',
                status: 'active',
                emailVerified: true,
                createdAt: like('2023-01-01T00:00:00Z'),
                updatedAt: like('2023-01-01T00:00:00Z'),
              },
              token: like('jwt-token'),
            },
            message: 'Login successful',
            errors: [],
            meta: {
              timestamp: like('2023-01-01T00:00:00Z'),
              requestId: like('550e8400-e29b-41d4-a716-446655440000'),
            },
          },
        },
      });

      const response = await axios.post(`${provider.mockService.baseUrl}/api/v1/auth/login`, {
        email: 'test@example.com',
        password: 'password123',
      });

      expect(response.status).toEqual(200);
      expect(response.data.success).toBe(true);
      expect(response.data.data.user.email).toEqual('test@example.com');
      expect(response.data.data.token).toBeDefined();
    });
  });
});
```

## 13. API Monitoring

### 13.1 Logging

Structured logging dengan **Logrus** atau **Zap**:

```go
// logger.go
package main

import (
  "net/http"
  "time"

  "github.com/sirupsen/logrus"
)

func loggingMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    
    // Log request
    logger.WithFields(logrus.Fields{
      "method": r.Method,
      "path":   r.URL.Path,
      "ip":     r.RemoteAddr,
      "agent":  r.UserAgent(),
    }).Info("Request started")
    
    // Create response recorder to capture status code
    recorder := httptest.NewRecorder()
    
    // Call next handler
    next.ServeHTTP(recorder, r)
    
    // Copy response to original writer
    for k, v := range recorder.Header() {
      w.Header()[k] = v
    }
    w.WriteHeader(recorder.Code)
    w.Write(recorder.Body.Bytes())
    
    // Log response
    duration := time.Since(start)
    logger.WithFields(logrus.Fields{
      "method":   r.Method,
      "path":     r.URL.Path,
      "status":   recorder.Code,
      "duration": duration,
    }).Info("Request completed")
  })
}
```

### 13.2 Metrics

Metrics collection dengan **Prometheus**:

```go
// metrics.go
package main

import (
  "strconv"
  "time"

  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promauto"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
  httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
    Name: "http_requests_total",
    Help: "Total number of HTTP requests",
  }, []string{"method", "endpoint", "status"})

  httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
    Name:    "http_request_duration_seconds",
    Help:    "Duration of HTTP requests",
    Buckets: prometheus.DefBuckets,
  }, []string{"method", "endpoint"})

  activeUsers = promauto.NewGauge(prometheus.GaugeOpts{
    Name: "active_users",
    Help: "Number of active users",
  })

  activeRooms = promauto.NewGauge(prometheus.GaugeOpts{
    Name: "active_rooms",
    Help: "Number of active rooms",
  })
)

func metricsMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    
    // Create response recorder to capture status code
    recorder := httptest.NewRecorder()
    
    // Call next handler
    next.ServeHTTP(recorder, r)
    
    // Copy response to original writer
    for k, v := range recorder.Header() {
      w.Header()[k] = v
    }
    w.WriteHeader(recorder.Code)
    w.Write(recorder.Body.Bytes())
    
    // Record metrics
    duration := time.Since(start).Seconds()
    status := strconv.Itoa(recorder.Code)
    
    httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
    httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
  })
}

func setupMetrics() {
  http.Handle("/metrics", promhttp.Handler())
}
```

### 13.3 Distributed Tracing

Distributed tracing dengan **Jaeger**:

```go
// tracing.go
package main

import (
  "net/http"

  "github.com/opentracing/opentracing-go"
  "github.com/uber/jaeger-client-go"
  jaegercfg "github.com/uber/jaeger-client-go/config"
)

func initTracer(serviceName string) (opentracing.Tracer, io.Closer, error) {
  cfg := jaegercfg.Configuration{
    ServiceName: serviceName,
    Sampler: &jaegercfg.SamplerConfig{
      Type:  jaeger.SamplerTypeConst,
      Param: 1,
    },
    Reporter: &jaegercfg.ReporterConfig{
      LogSpans: true,
    },
  }
  
  return cfg.NewTracer()
}

func tracingMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Extract span context from headers
    wireContext, err := opentracing.GlobalTracer().Extract(
      opentracing.HTTPHeaders,
      opentracing.HTTPHeadersCarrier(r.Header),
    )
    
    var span opentracing.Span
    if err != nil {
      // No span context in headers, create root span
      span = opentracing.StartSpan(r.URL.Path)
    } else {
      // Create child span
      span = opentracing.StartSpan(r.URL.Path, opentracing.ChildOf(wireContext))
    }
    defer span.Finish()
    
    // Inject span context into request context
    ctx := opentracing.ContextWithSpan(r.Context(), span)
    
    // Call next handler with context
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}
```

## 14. Kesimpulan

Desain API endpoints untuk aplikasi WebRTC meeting ini dirancang dengan mempertimbangkan:

1. **RESTful design** dengan resource-oriented URLs dan proper HTTP methods
2. **Consistent response format** untuk semua endpoints
3. **Comprehensive error handling** dengan status codes yang tepat
4. **Security features** seperti JWT authentication, rate limiting, dan input validation
5. **Real-time communication** melalui WebSocket untuk WebRTC signaling
6. **Documentation** dengan OpenAPI/Swagger specification
7. **Testing strategy** dengan unit, integration, dan contract testing
8. **Monitoring dan logging** untuk observability

Dengan desain API ini, frontend (Preact) dapat berkomunikasi dengan backend (Golang) secara efisien dan aman untuk mendukung semua fitur aplikasi WebRTC meeting.