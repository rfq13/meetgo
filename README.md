# WebRTC Meeting Backend

A comprehensive WebRTC meeting backend server built with Golang, providing real-time video conferencing capabilities with WebSocket signaling, PostgreSQL database, and Janus WebRTC server integration.

## 🚀 Features

- **Real-time Video Conferencing**: WebRTC-based video/audio communication
- **WebSocket Signaling**: Real-time signaling for WebRTC connections
- **User Management**: User registration, authentication, and profile management
- **Room Management**: Create, join, and manage meeting rooms
- **PostgreSQL Database**: Robust data persistence with GORM ORM
- **Janus WebRTC Server**: Professional WebRTC media server integration
- **RESTful API**: Complete REST API for all operations
- **JWT Authentication**: Secure token-based authentication
- **Docker Support**: Containerized deployment with Docker Compose
- **Comprehensive Testing**: Full test suite for all components

## 📋 Prerequisites

- **Go 1.21+**: Go programming language
- **PostgreSQL 13+**: Database server
- **Docker & Docker Compose**: For containerized deployment
- **Node.js 18+**: For running WebSocket client tests (optional)
- **Make**: For build automation

## 🛠️ Installation

### 1. Clone the Repository

```bash
git clone <repository-url>
cd webrtc
```

### 2. Environment Configuration

Copy the environment template and configure:

```bash
cp .env.example .env
```

Edit `.env` file with your configuration:

```env
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8080
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s
SERVER_IDLE_TIMEOUT=60s

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=webrtc_meeting
DB_SSLMODE=disable

# WebSocket Server Configuration
WS_HOST=localhost
WS_PORT=8081

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key
JWT_REFRESH_SECRET=your-super-secret-refresh-key
JWT_EXPIRATION_TIME=24h
JWT_REFRESH_TIME=168h

# Janus WebRTC Server Configuration
JANUS_BASE_URL=http://localhost:8088/janus
JANUS_ADMIN_URL=http://localhost:8088/admin
JANUS_WS_URL=ws://localhost:8188
JANUS_API_SECRET=janusrocks
JANUS_ADMIN_SECRET=janusrocks

# Redis Configuration (optional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Email Configuration (optional)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
EMAIL_FROM=noreply@webrtc-meeting.com

# Logger Configuration
LOG_LEVEL=info
LOG_FORMAT=json
```

### 3. Database Setup

#### Option A: Local PostgreSQL

```bash
# Create database
createdb webrtc_meeting

# Run migrations (handled automatically by the application)
```

#### Option B: Docker PostgreSQL

```bash
docker run -d \
  --name postgres-webrtc \
  -e POSTGRES_DB=webrtc_meeting \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=your_password \
  -p 5432:5432 \
  postgres:13
```

### 4. Janus WebRTC Server Setup

#### Option A: Docker (Recommended)

```bash
docker run -d \
  --name janus \
  -p 8088:8088 \
  -p 8188:8188 \
  -p 10000-20000:10000-20000/udp \
  januswebrtc/janus:latest
```

#### Option B: Local Installation

Follow the official Janus WebRTC Server installation guide:
[https://janus.conf.meetecho.com/docs/](https://janus.conf.meetecho.com/docs/)

## 🔨 Building and Running

### Using Make Commands (Recommended)

```bash
# Build all binaries
make build

# Run API server
make run-api

# Run WebSocket server
make run-websocket

# Run tests
make test

# Run validation
make validate
```

### Manual Build and Run

```bash
# Navigate to backend directory
cd backend

# Download dependencies
go mod download
go mod tidy

# Build binaries
go build -o bin/api-server ./cmd/api
go build -o bin/websocket-server ./cmd/websocket

# Run API server
./bin/api-server

# Run WebSocket server (in separate terminal)
./bin/websocket-server
```

### Using Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## 🧪 Testing and Validation

The project includes comprehensive test scripts for validation:

### Available Test Scripts

1. **Compilation Test**: Validates Go code compilation and dependencies
2. **Database Test**: Tests database connection and migrations
3. **API Test**: Validates REST API endpoints
4. **WebSocket Test**: Tests WebSocket server functionality
5. **Janus Test**: Tests Janus WebRTC server integration

### Running Individual Tests

```bash
# Make scripts executable
chmod +x scripts/*.sh

# Run compilation test
./scripts/test-compilation.sh

# Run database test
./scripts/test-database.sh

# Run API test
./scripts/test-api.sh

# Run WebSocket test
./scripts/test-websocket.sh

# Run Janus test
./scripts/test-janus.sh
```

### Running All Validations

```bash
# Run all validation scripts
./scripts/validate-all.sh

# Or using Make
make validate
```

### Test Reports

Each test script generates a detailed report file:
- `build-report-YYYYMMDD-HHMMSS.txt`
- `database-test-report-YYYYMMDD-HHMMSS.txt`
- `api-test-report-YYYYMMDD-HHMMSS.txt`
- `websocket-test-report-YYYYMMDD-HHMMSS.txt`
- `janus-test-report-YYYYMMDD-HHMMSS.txt`

## 📚 API Documentation

### Base URLs

- **API Server**: `http://localhost:8080`
- **WebSocket Server**: `ws://localhost:8081`

### Authentication Endpoints

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "username",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Refresh Token
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "refresh_token_here"
}
```

#### Logout
```http
POST /api/v1/auth/logout
Authorization: Bearer <access_token>
```

### User Endpoints

#### Get User Profile
```http
GET /api/v1/users/profile
Authorization: Bearer <access_token>
```

#### Update User Profile
```http
PUT /api/v1/users/profile
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890"
}
```

### Room Endpoints

#### Create Room
```http
POST /api/v1/rooms
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "Meeting Room",
  "description": "Team meeting",
  "max_users": 10,
  "type": "meeting",
  "is_public": true
}
```

#### Get Rooms
```http
GET /api/v1/rooms
Authorization: Bearer <access_token>
```

#### Get Room Details
```http
GET /api/v1/rooms/{room_id}
Authorization: Bearer <access_token>
```

### Health Check Endpoints

```http
GET /health
GET /ping
GET /ready
GET /live
```

### WebSocket Endpoints

#### Connect to WebSocket
```
ws://localhost:8081/ws?userId={user_id}
```

#### Authenticated WebSocket Connection
```
ws://localhost:8081/ws/auth?token={access_token}
```

#### WebSocket Message Format

```json
{
  "type": "join_room|leave_room|signal|chat_message",
  "roomId": "room_id",
  "data": {
    "message": "message content"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 🏗️ Project Structure

```
webrtc/
├── backend/                    # Go backend application
│   ├── cmd/                   # Application entry points
│   │   ├── api/              # API server main
│   │   └── websocket/        # WebSocket server main
│   ├── internal/             # Internal application code
│   │   ├── api/             # API handlers and routes
│   │   ├── auth/            # Authentication logic
│   │   ├── config/          # Configuration management
│   │   ├── database/        # Database connection and migrations
│   │   ├── models/          # Data models
│   │   ├── room/            # Room management
│   │   ├── user/            # User management
│   │   ├── webrtc/          # WebRTC and Janus integration
│   │   └── websocket/       # WebSocket handling
│   ├── pkg/                 # Public packages
│   │   └── logger/          # Logging utilities
│   ├── models/              # Database models
│   ├── configs/             # Configuration files
│   ├── go.mod               # Go modules
│   └── go.sum               # Go dependencies
├── frontend/                 # Frontend application (React)
├── janus-server/            # Janus configuration
│   ├── config/              # Janus config files
│   └── Dockerfile           # Janus Docker setup
├── scripts/                 # Test and validation scripts
│   ├── test-compilation.sh  # Compilation test
│   ├── test-database.sh     # Database test
│   ├── test-api.sh          # API test
│   ├── test-websocket.sh    # WebSocket test
│   ├── test-janus.sh        # Janus test
│   └── validate-all.sh      # Run all tests
├── docker-compose.yml       # Docker Compose configuration
├── Makefile                # Build automation
├── .env.example            # Environment template
└── README.md               # This file
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_HOST` | API server host | `localhost` |
| `SERVER_PORT` | API server port | `8080` |
| `WS_HOST` | WebSocket server host | `localhost` |
| `WS_PORT` | WebSocket server port | `8081` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database username | `postgres` |
| `DB_PASSWORD` | Database password | - |
| `DB_NAME` | Database name | `webrtc_meeting` |
| `JWT_SECRET` | JWT signing secret | - |
| `JANUS_BASE_URL` | Janus HTTP API URL | `http://localhost:8088/janus` |
| `JANUS_WS_URL` | Janus WebSocket URL | `ws://localhost:8188` |

### Database Configuration

The application uses PostgreSQL with the following main tables:

- `users` - User accounts and profiles
- `user_sessions` - User authentication sessions
- `rooms` - Meeting rooms
- `room_participants` - Room participants
- `room_messages` - Chat messages
- `meeting_history` - Meeting records

## 🚀 Deployment

### Docker Deployment

```bash
# Build and start all services
docker-compose up -d

# Scale services
docker-compose up -d --scale api=2 --scale websocket=2

# View logs
docker-compose logs -f api
docker-compose logs -f websocket
```

### Production Deployment

1. **Environment Setup**: Configure production environment variables
2. **Database**: Set up PostgreSQL cluster
3. **Load Balancer**: Configure load balancer for API and WebSocket servers
4. **SSL/TLS**: Configure SSL certificates
5. **Monitoring**: Set up monitoring and logging

### Health Checks

The application provides multiple health check endpoints:

- `/health` - Overall system health
- `/ping` - Simple ping response
- `/ready` - Readiness probe
- `/live` - Liveness probe

## 🐛 Troubleshooting

### Common Issues

#### 1. Database Connection Failed

```bash
# Check PostgreSQL status
docker ps | grep postgres

# Check database logs
docker logs postgres-webrtc

# Test connection manually
psql -h localhost -U postgres -d webrtc_meeting
```

#### 2. Janus Server Not Responding

```bash
# Check Janus container
docker ps | grep janus

# Check Janus logs
docker logs janus

# Test Janus API
curl http://localhost:8088/janus/info
```

#### 3. WebSocket Connection Issues

```bash
# Check WebSocket server
curl http://localhost:8081/health

# Test WebSocket connection
wscat -c ws://localhost:8081/ws?userId=test
```

#### 4. Build Errors

```bash
# Clean Go modules
go clean -modcache
go mod download

# Rebuild
make build
```

### Debug Mode

Enable debug logging by setting:

```env
LOG_LEVEL=debug
LOG_FORMAT=text
```

## 📈 Monitoring and Logging

### Logging

The application uses structured logging with logrus:

- **JSON Format**: For production environments
- **Text Format**: For development environments
- **Log Levels**: debug, info, warn, error, fatal

### Metrics

The application exposes metrics for monitoring:

- Connection counts
- Request/response times
- Error rates
- Database performance

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Run validation: `make validate`
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support and questions:

1. Check the troubleshooting section
2. Review test reports for validation results
3. Check application logs
4. Create an issue in the repository

## 🔄 Version History

- **v1.0.0** - Initial release with core WebRTC functionality
  - User authentication and management
  - Room creation and management
  - WebSocket signaling
  - Janus WebRTC integration
  - Comprehensive test suite

## 📚 Additional Documentation

- [API Documentation](docs/api.md)
- [Database Schema](docs/database.md)
- [Deployment Guide](docs/deployment.md)
- [Architecture Overview](docs/architecture.md)

---

**Built with ❤️ using Go, WebRTC, and Janus**