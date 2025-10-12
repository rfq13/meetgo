# WebRTC Meeting - Docker Configuration

This document explains how to set up and run the WebRTC Meeting application using Docker and Docker Compose.

## Prerequisites

Before you begin, ensure you have the following installed:

- [Docker](https://docs.docker.com/get-docker/) (version 20.10 or later)
- [Docker Compose](https://docs.docker.com/compose/install/) (version 2.0 or later)
- [Make](https://www.gnu.org/software/make/) (optional, for convenient commands)

## Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd webrtc-meeting
```

### 2. Setup Environment

```bash
# Copy environment variables template
cp .env.example .env

# Edit the .env file with your configuration
nano .env
```

### 3. Start the Application

#### Option A: Using Make (Recommended)

```bash
# Start development environment
make dev

# View logs
make docker-logs

# Stop services
make docker-down
```

#### Option B: Using Docker Compose Directly

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### 4. Access the Services

Once started, you can access the following services:

- **Frontend**: http://localhost:3000
- **API Server**: http://localhost:8080
- **WebSocket Server**: ws://localhost:8081
- **Janus WebRTC Server**: http://localhost:8088/janus
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

## Architecture

The application consists of the following services:

### Core Services

1. **PostgreSQL Database** (`postgres`)
   - Port: 5432
   - Database: `webrtc_meeting`
   - User: `webrtc_user`

2. **Redis Cache** (`redis`)
   - Port: 6379
   - Password: configured in `.env`

3. **Backend API Server** (`api`)
   - Port: 8080
   - Handles REST API requests
   - Built with Golang and Gin framework

4. **WebSocket Server** (`websocket`)
   - Port: 8081
   - Handles real-time communication
   - WebRTC signaling

5. **Janus WebRTC Server** (`janus`)
   - HTTP Port: 8088
   - HTTPS Port: 7889
   - WebSocket Port: 8189
   - RTP/UDP Ports: 10000-20000

### Optional Services

6. **Frontend** (`frontend`)
   - Port: 3000
   - Built with Preact/React
   - Development server only

7. **Nginx Reverse Proxy** (`nginx`)
   - HTTP Port: 80
   - HTTPS Port: 443
   - Production use only

## Configuration

### Environment Variables

Key environment variables in `.env`:

```bash
# Database
POSTGRES_DB=webrtc_meeting
POSTGRES_USER=webrtc_user
POSTGRES_PASSWORD=your_secure_password

# Redis
REDIS_PASSWORD=your_redis_password

# JWT
JWT_SECRET=your_super_secret_jwt_key_min_32_chars

# Janus
JANUS_API_SECRET=janusrocks
JANUS_ADMIN_SECRET=janusrocksadmin

# Ports
API_PORT=8080
WEBSOCKET_PORT=8081
FRONTEND_PORT=3000
```

### Janus Configuration

Janus WebRTC Server configuration files:

- `janus-server/config/janus.jcfg` - Main Janus configuration
- `janus-server/config/janus.plugin.videoroom.jcfg` - VideoRoom plugin configuration

## Make Commands

The Makefile provides convenient commands for development and deployment:

### Setup Commands

```bash
make install          # Install dependencies
make setup            # Setup project (copy .env)
make help             # Show all available commands
```

### Development Commands

```bash
make dev              # Start development environment
make dev-backend      # Start only backend services
make dev-frontend     # Start only frontend service
make docker-logs      # Show all logs
make docker-down      # Stop all services
```

### Build Commands

```bash
make build            # Build all applications
make build-backend    # Build backend only
make build-frontend   # Build frontend only
make docker-build     # Build Docker images
```

### Testing Commands

```bash
make test             # Run all tests
make test-backend     # Run backend tests
make test-frontend    # Run frontend tests
make test-coverage    # Run tests with coverage
```

### Quality Commands

```bash
make lint             # Run linting
make format           # Format code
make clean            # Clean build artifacts
```

### Database Commands

```bash
make db-migrate       # Run database migrations
make db-reset         # Reset database
make db-backup        # Backup database
make db-restore       # Restore database (requires BACKUP_FILE parameter)
```

### Monitoring Commands

```bash
make health           # Check service health
make status           # Show service status
make logs-api         # Show API logs
make logs-websocket   # Show WebSocket logs
make logs-janus       # Show Janus logs
```

### Utility Commands

```bash
make shell-api        # Open shell in API container
make shell-postgres   # Open shell in PostgreSQL
make shell-redis      # Open shell in Redis
make version          # Show version information
make info             # Show project information
```

## Development Workflow

### 1. Initial Setup

```bash
# Install dependencies
make install

# Setup environment
make setup

# Start development environment
make dev
```

### 2. Development

```bash
# View logs
make docker-logs

# Access service shells
make shell-api
make shell-postgres

# Run tests
make test

# Format code
make format
```

### 3. Building for Production

```bash
# Build Docker images
make docker-build

# Start production environment
make prod
```

## Docker Compose Profiles

The `docker-compose.yml` uses profiles to manage different deployment scenarios:

### Development Profile

```bash
# Start with frontend (default)
docker-compose --profile frontend up -d
```

### Production Profile

```bash
# Start with nginx reverse proxy
docker-compose --profile production up -d
```

## Network Configuration

All services communicate through a dedicated Docker network:

- **Network Name**: `webrtc-network`
- **Subnet**: `172.20.0.0/16`
- **Driver**: `bridge`

## Volume Management

Persistent data is stored in Docker volumes:

- `postgres_data`: PostgreSQL data
- `redis_data`: Redis data

Additional bind mounts:

- `./backend/logs:/app/logs`: Backend logs
- `./janus-server/logs:/var/log/janus`: Janus logs
- `./janus-server/recordings:/opt/janus/share/janus/recordings`: Recordings

## Health Checks

All services include health checks:

- **API Server**: `GET /health`
- **WebSocket Server**: `GET /health`
- **Janus Server**: `GET /janus/info`
- **PostgreSQL**: `pg_isready`
- **Redis**: `ping`

## Troubleshooting

### Common Issues

1. **Port Conflicts**
   ```bash
   # Check which ports are in use
   netstat -tulpn | grep :8080
   
   # Change ports in .env file
   API_PORT=8081
   ```

2. **Permission Issues**
   ```bash
   # Fix Docker permissions
   sudo usermod -aG docker $USER
   
   # Restart Docker service
   sudo systemctl restart docker
   ```

3. **Database Connection Issues**
   ```bash
   # Reset database
   make db-reset
   
   # Check database logs
   make logs-postgres
   ```

4. **Janus WebRTC Issues**
   ```bash
   # Check Janus logs
   make logs-janus
   
   # Test Janus API
   curl http://localhost:8088/janus/info
   ```

### Debug Commands

```bash
# Check service status
make status

# Check service health
make health

# View detailed logs
docker-compose logs -f [service-name]

# Access service shell
make shell-api
make shell-websocket
make shell-janus
```

## Security Considerations

### Production Deployment

1. **Change Default Secrets**
   ```bash
   # Update these in .env
   JWT_SECRET=your_secure_random_string_min_32_chars
   JANUS_API_SECRET=your_janus_api_secret
   JANUS_ADMIN_SECRET=your_janus_admin_secret
   POSTGRES_PASSWORD=your_secure_db_password
   REDIS_PASSWORD=your_secure_redis_password
   ```

2. **Use HTTPS**
   ```bash
   # Enable HTTPS in production
   make prod
   ```

3. **Network Security**
   - Use firewall to restrict access to database ports
   - Configure CORS properly in production
   - Use VPN for remote access

4. **Regular Updates**
   ```bash
   # Update Docker images
   docker-compose pull
   
   # Restart services
   docker-compose up -d
   ```

## Performance Optimization

### Resource Limits

Adjust resource limits in `docker-compose.yml`:

```yaml
services:
  api:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

### Monitoring

Enable monitoring in production:

```bash
# Set environment variables
METRICS_ENABLED=true
METRICS_PORT=9090

# Access metrics
curl http://localhost:9090/metrics
```

## Backup and Recovery

### Database Backup

```bash
# Create backup
make db-backup

# Manual backup
docker-compose exec postgres pg_dump -U webrtc_user webrtc_meeting > backup.sql

# Restore backup
make db-restore BACKUP_FILE=backup.sql
```

### Configuration Backup

```bash
# Backup configuration files
tar -czf config-backup.tar.gz .env janus-server/config/

# Restore configuration
tar -xzf config-backup.tar.gz
```

## Contributing

When contributing to the Docker configuration:

1. Test changes with `docker-compose config`
2. Update this README if adding new services
3. Ensure all health checks work properly
4. Test both development and production profiles

## Support

For issues related to:

- **Docker Configuration**: Check this README and troubleshooting section
- **Application Issues**: Check the main README.md
- **Janus WebRTC**: Check [Janus Documentation](https://janus.conf.meetecho.com/docs/)
- **Docker Issues**: Check [Docker Documentation](https://docs.docker.com/)

## License

This Docker configuration is part of the WebRTC Meeting application. See the main LICENSE file for details.