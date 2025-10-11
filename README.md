<div align="center">
# SoundTube - Audio Sharing Platform 🔊

![Status](https://img.shields.io/badge/Status-Active-success)
![License](https://img.shields.io/badge/License-Proprietary-red)
![Version](https://img.shields.io/badge/Version-0.16B-blue)
![Platform](https://img.shields.io/badge/Platform-Web-informational)

A full-featured audio sharing and streaming platform built with Go, Gin, and PostgreSQL. Users can upload, share, and interact with audio content including likes, dislikes, and comments.

</div>

### Lead Architect
**Finimen Sniper** - 📧 finimensniper@gmail.com

## 🚀 Features

### Core Functionality
- **User Authentication** - JWT-based registration/login with email verification
- **Audio Management** - Upload, stream, and manage audio files
- **Social Features** - Like/dislike system and comments
- **File Handling** - Secure audio file upload and storage

### Technical Features
- **Rate Limiting** - IP-based request throttling
- **Caching** - Redis for performance optimization
- **Tracing** - OpenTelemetry integration for observability
- **Security** - Middleware for CORS, JWT validation, and secure headers
- **Health Checks** - Comprehensive service monitoring

## 🛠 Tech Stack

### Backend
- **Go 1.21+** - Primary programming language
- **Gin** - HTTP web framework
- **PostgreSQL** - Primary database
- **Redis** - Caching and token blacklisting
- **JWT** - Authentication tokens

### Infrastructure
- **Docker** - Containerization
- **OpenTelemetry** - Distributed tracing
- **Viper** - Configuration management

## 📁 Project Structure

```
soundtube/
├── cmd/
│   └── di/                 # Dependency injection container
├── internal/
│   ├── domain/             # Domain models and interfaces
│   │   ├── auth/           # Authentication domain
│   │   ├── sound/          # Sound domain
│   │   └── reactions/      # Reactions domain
│   ├── handlers/           # HTTP request handlers
│   ├── services/           # Business logic layer
│   └── repositories/       # Data access layer
├── pkg/
│   ├── config/             # Configuration management
│   ├── middleware/         # HTTP middleware
│   └── utils/              # Shared utilities
├── configs/                # Configuration files
├── static/                 # Static files and uploads
└── migrations/             # Database migrations
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 14+
- Redis 6+
- Docker (optional)

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/your-username/soundtube.git
cd soundtube
```

2. **Set up configuration**
```bash
cp configs/dev.example.yaml configs/dev.yaml
# Edit configs/dev.yaml with your settings
```

3. **Configure environment**
```yaml
# configs/dev.yaml
database:
  host: "localhost"
  port: "5432"
  user: "soundtube_user"
  password: "your_password"
  dbname: "soundtube"

redis:
  addr: "localhost:6379"
  password: ""
  db: 0

server:
  host: "localhost"
  port: ":8080"
```

4. **Run database migrations**
```sql
-- The application automatically runs embedded migrations
-- Manual execution if needed:
psql -d soundtube -f migrations/001_initial_schema.sql
```

5. **Start the application**
```bash
go run cmd/main.go
```

### Docker Setup

```bash
# Using Docker Compose
docker-compose up -d

# Or build manually
docker build -t soundtube .
docker run -p 8080:8080 soundtube
```

## 🔧 API Documentation

### Authentication Endpoints

<div align="center">
  
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | User login |
| POST | `/api/auth/logout` | User logout |
| GET | `/api/auth/verify-email` | Verify email address |

### Sounds Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/sounds` | Get all sounds |
| POST | `/api/sounds` | Create sound record |
| POST | `/api/sounds/upload` | Upload audio file |
| PATCH | `/api/sounds/{id}` | Update sound |
| DELETE | `/api/sounds/{id}` | Delete sound |

### Reactions Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| PUT | `/api/sounds/{id}/reactions` | Add reaction to sound |
| DELETE | `/api/sounds/{id}/reactions` | Remove reaction from sound |
| GET | `/api/sounds/{id}/reactions` | Get sound reactions |

</div>

### Example Requests

**User Registration:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

**Upload Sound:**
```bash
curl -X POST http://localhost:8080/api/sounds/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@/path/to/audio.mp3" \
  -F "name=My Awesome Sound"
```

## 🔒 Security Features

- **JWT Authentication** with configurable expiration
- **Password Hashing** using bcrypt
- **Rate Limiting** per IP address
- **CORS Protection**
- **Secure Headers** middleware
- **Token Blacklisting** for logout functionality

## 📊 Monitoring & Observability

### Health Endpoints
- `GET /health` - Comprehensive health check
- `GET /ready` - Readiness probe
- `GET /live` - Liveness probe

### Tracing
The application supports OpenTelemetry tracing with Jaeger. Enable in config:

```yaml
traycing:
  enabled: true
  service_name: "soundtube-api"
  endpoint: "http://localhost:14268/api/traces"
```

## 🗄 Database Schema

### Key Tables
- `users` - User accounts and profiles
- `sounds` - Audio metadata and file information
- `sound_reactions` - Like/dislike counts
- `sound_participants` - User reaction tracking
- `comments` - User comments on sounds

## 🧪 Testing

```bash
# Run unit tests
go test ./...

# Run integration tests
go test -tags=integration ./...

# Test with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🔧 Configuration

### Environment Variables
All configuration is managed through YAML files, but can be overridden with environment variables:

```bash
export SOUNDTUBE_DATABASE_HOST=localhost
export SOUNDTUBE_REDIS_ADDR=localhost:6379
export SOUNDTUBE_SERVER_PORT=:8080
```

### Key Configuration Sections
- **Database** - Connection pooling and timeouts
- **Redis** - Cache and session storage
- **JWT** - Token signing and expiration
- **Rate Limiting** - Request thresholds
- **Email** - SMTP configuration for verification

## 🚀 Deployment

### Production Build
```bash
# Build optimized binary
GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/soundtube cmd/main.go

# Run in production mode
export GIN_MODE=release
./bin/soundtube
```

### Kubernetes (Example)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: soundtube-api
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: soundtube
        image: soundtube:latest
        ports:
        - containerPort: 8080
        env:
        - name: GIN_MODE
          value: "release"
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow Go best practices and effective Go guidelines
- Write tests for new functionality
- Update documentation for API changes
- Use conventional commit messages

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/Finimen/Soundtube/blob/main/License.md) file for details.

## 🆘 Support

- 📧 Email: support@soundtube.com
- 🐛 Issues: [GitHub Issues](https://github.com/your-username/soundtube/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/your-username/soundtube/discussions)

## 🙏 Acknowledgments

- Gin Web Framework community
- PostgreSQL and Redis communities
- OpenTelemetry for observability tools

---

**SoundTube** - Share your sound with the world! 🎵

