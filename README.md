# Smart Log Viewer 🚀

A real-time log streaming application built with **Go** (WebSocket server) and **React TypeScript** (client), featuring real-time log streaming, pause/resume functionality, and comprehensive DevOps automation.

## 🌟 Features

- **Real-time Log Streaming**: WebSocket-based live log streaming
- **Interactive UI**: React-based interface with log level filtering
- **Pause/Resume**: Control log flow without losing connection
- **Multi-level Logs**: INFO, WARN, ERROR with color coding
- **Auto-reconnection**: Robust WebSocket connection management
- **Responsive Design**: Modern, mobile-friendly interface

## 🏗️ Architecture

```
┌─────────────────┐    WebSocket    ┌─────────────────┐
│   React Client  │ ←──────────────→ │   Go Server     │
│   (Port 3000)   │                 │   (Port 8080)   │
└─────────────────┘                 └─────────────────┘
         │                                   │
         │                                   │
         ▼                                   ▼
┌─────────────────┐                 ┌─────────────────┐
│   Nginx Proxy   │                 │  Mock Log Gen   │
│   (Port 80/443) │                 │   (Every 1s)    │
└─────────────────┘                 └─────────────────┘
```

## 🚀 Quick Start

### Prerequisites
- **Docker** and **Docker Compose**
- **Go 1.21+** (for local development)
- **Node.js 18+** (for local development)

### Option 1: Docker Compose (Recommended)
```bash
# Clone the repository
git clone <your-repo-url>
cd smart-log-viewer

# Start the entire stack
docker-compose up -d

# Access the application
# Client: http://localhost:3000
# Server: http://localhost:8080
# WebSocket: ws://localhost:8080/ws
```

### Option 2: Development Mode
```bash
# Start development services with hot reload
docker-compose --profile development up -d

# Access development services
# Client Dev: http://localhost:3001
# Server Dev: http://localhost:8081
```

### Option 3: Local Development
```bash
# Terminal 1: Start Go server
cd server
go run ./cmd/server

# Terminal 2: Start React client
cd client
npm install
npm run dev
```

## 🔧 DevOps Components

### 1. Git Pre-commit Hook

**Purpose**: Automatically validates code quality before each commit.

**What it does**:
- ✅ Formats Go code with `go fmt`
- ✅ Lints Go code with `golangci-lint`
- ✅ Runs Go tests
- ✅ Builds Go binary
- ✅ Lints React/TypeScript code
- ✅ Runs client tests
- ✅ Builds client
- ✅ Validates Docker Compose configuration

**Location**: `.git/hooks/pre-commit`

**Usage**: Automatically runs on `git commit`. To bypass: `git commit --no-verify`

### 2. GitHub Actions CI/CD Pipeline

**Purpose**: Automatically builds, tests, and deploys code on every push/PR.

**What it does**:
- 🐹 **Test Server**: Go linting, testing, building
- ⚛️ **Test Client**: TypeScript linting, testing, building
- 🐳 **Build Docker**: Creates container images
- 🔗 **Integration Tests**: Tests services together
- 🔒 **Security Scan**: Vulnerability scanning with Trivy
- 🚀 **Deploy Staging**: Automatic staging deployment
- 🚀 **Deploy Production**: Manual production deployment

**Location**: `.github/workflows/ci.yml`

**Triggers**: 
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop`

### 3. Docker Compose Orchestration

**Purpose**: Single command to start entire application stack.

**What it does**:
- 🐳 **Multi-service**: Client, server, nginx
- 🔄 **Profiles**: Development vs production modes
- 🌐 **Networking**: Isolated network with custom subnet
- 📊 **Health Checks**: Automatic service monitoring
- 🔒 **Volumes**: Persistent data storage
- 🚀 **Restart Policy**: Automatic service recovery

**Location**: `docker-compose.yml`

**Profiles**:
- `production`: Optimized for production use
- `development`: Hot reload and debugging

## 📁 Project Structure

```
smart-log-viewer/
├── .github/
│   └── workflows/
│       └── ci.yml                 # GitHub Actions CI/CD
├── .git/
│   └── hooks/
│       └── pre-commit            # Git pre-commit hook
├── client/                        # React TypeScript client
│   ├── src/
│   ├── Dockerfile
│   ├── package.json
│   └── tsconfig.json
├── server/                        # Go WebSocket server
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── model/
│   │   ├── websocket/
│   │   └── loggenerator/
│   ├── Dockerfile
│   └── go.mod
├── docker-compose.yml             # Multi-service orchestration
├── README.md                      # This file
└── .gitignore
```

## 🐳 Docker Services

### Production Services
- **server**: Go WebSocket server (port 8080)
- **client**: React application (port 3000)
- **nginx**: Reverse proxy (port 80/443)

### Development Services
- **server-dev**: Go server with hot reload (port 8081)
- **client-dev**: React with hot reload (port 3001)

## 🔍 Monitoring & Health Checks

### Health Check Endpoints
- **Server**: `http://localhost:8080/` (status page)
- **Client**: `http://localhost:3000/` (React app)
- **WebSocket**: `ws://localhost:8080/ws` (real-time logs)

### Docker Health Checks
- **Server**: HTTP health check every 30s
- **Client**: HTTP health check every 30s


## 🚀 Deployment

### Staging (Automatic)
- Triggers on push to `main` branch
- Runs after all tests pass
- Deploys to staging environment

### Production (Manual)
- Requires manual approval
- Runs after staging deployment
- Deploys to production environment

## 🛠️ Development

### Adding New Features
1. Create feature branch: `git checkout -b feature/new-feature`
2. Make changes and test locally
3. Pre-commit hook will validate code quality
4. Push and create PR
5. CI pipeline will test and build
6. Merge after approval

### Testing
```bash
# Run Go tests
cd server
go test ./...

# Run client tests
cd client
npm test

# Run integration tests
docker-compose up -d
# Tests run automatically in CI
```

### Linting
```bash
# Go linting
cd server
golangci-lint run ./...

# Client linting
cd client
npm run lint
```

## 🔧 Configuration

### Environment Variables
- `GO_ENV`: Go environment (development/production)
- `NODE_ENV`: Node.js environment (development/production)
- `REACT_APP_WS_URL`: WebSocket server URL
- `PORT`: Server port (default: 8080)

### Docker Compose Profiles
- **Default**: Production mode
- **Development**: Hot reload and debugging
- **Custom**: Mix and match services

## 📊 Performance

### Optimizations
- **Go**: Efficient WebSocket handling with goroutines
- **React**: Virtual DOM and optimized rendering
- **Docker**: Multi-stage builds and layer caching
- **CI/CD**: Parallel job execution and caching

### Scalability
- **Horizontal**: Add more server instances
- **Vertical**: Increase container resources
- **Load Balancing**: Nginx reverse proxy


## 🚨 Troubleshooting

### Common Issues

#### WebSocket Connection Failed
```bash
# Check server status
curl http://localhost:8080/

# Check server logs
docker-compose logs server

# Verify ports
netstat -tulpn | grep 8080
```

#### Client Build Failed
```bash
# Clear node modules
cd client
rm -rf node_modules package-lock.json
npm install

# Check Node.js version
node --version  # Should be 18+
```

#### Docker Compose Issues
```bash
# Reset everything
docker-compose down -v
docker system prune -f
docker-compose up -d

# Check service status
docker-compose ps
```

### Logs
```bash
# View all logs
docker-compose logs

# View specific service
docker-compose logs server
docker-compose logs client

# Follow logs
docker-compose logs -f server
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch
3. Make changes
4. Run pre-commit checks
5. Create pull request
6. Wait for CI pipeline
7. Get approval and merge

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- **Go**: For efficient WebSocket server
- **React**: For modern UI framework
- **Docker**: For containerization
- **GitHub Actions**: For CI/CD automation

---

**Happy Log Viewing! 🎉**
