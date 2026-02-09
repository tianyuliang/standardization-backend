# standardization-backend

> A microservices-based standardization management backend system built with Go-Zero framework.

## Overview

This is a standardization management backend service that provides comprehensive data standardization capabilities including:

- **Rule Management** (`rule-api`) - Encoding rule management
- **Catalog Management** (`catalog-api`) - Standard catalog hierarchy management
- **Standard File Management** (`stdfile-api`) - Standard document/file management
- **Dictionary Management** (`dict-api`) - Code dictionary management
- **Data Element Management** (`dataelement-api`) - Data element lifecycle management
- **Task Management** (`task-api`) - Standard task management

## Tech Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| Language | Go | 1.24+ |
| Framework | Go-Zero | v1.9+ |
| Database | MySQL | 8.0 |
| ORM | GORM + SQLx | latest |
| Message Queue | Kafka | 3.0 |
| API Documentation | Swagger/OpenAPI | 3.0 |

## Project Structure

```
standardization-backend/
├── api/                          # API Gateway Service
│   ├── doc/                      # API definitions (.api files)
│   │   ├── api.api              # Main API entry
│   │   ├── base.api             # Common types
│   │   ├── rule/                # Rule module API
│   │   ├── catalog/             # Catalog module API
│   │   ├── stdfile/             # Standard file module API
│   │   ├── dict/                # Dictionary module API
│   │   ├── dataelement/         # Data element module API
│   │   └── task/                # Task module API
│   ├── internal/
│   │   ├── handler/             # HTTP handlers
│   │   ├── logic/               # Business logic
│   │   ├── svc/                 # Service context
│   │   ├── types/               # Type definitions
│   │   └── errorx/              # Error handling
│   └── etc/                     # Configuration files
├── model/                        # Data models
├── migrations/                   # Database migrations
├── specs/                        # SDD specification documents
│   ├── 1-rule-api/
│   ├── 2-catalog-api/
│   ├── 3-std-file-api/
│   ├── 4-dict-api/
│   ├── 5-std-task-api/
│   └── 6-std-dataelement-api/
├── deploy/                       # Deployment configurations
│   ├── docker/
│   └── k8s/
├── Makefile                      # Build automation
├── CLAUDE.md                     # Project guidelines
└── go.mod                        # Go module definition
```

## Quick Start

### Prerequisites

- Go 1.24 or higher
- MySQL 8.0
- Kafka 3.0 (optional, for async messaging)
- `goctl` CLI tool

### Installation

```bash
# Clone the repository
git clone https://github.com/tianyuliang/standardization-backend.git
cd standardization-backend

# Install dependencies
make deps

# Generate API code
make api

# Run the service
make run
```

The API will be available at `http://localhost:8888`

### Using Docker

```bash
# Build Docker image
make docker-build

# Run container
make docker-run
```

### Using Kubernetes

```bash
# Deploy to dev environment
make k8s-deploy-dev

# Deploy to prod environment
make k8s-deploy-prod

# Check status
make k8s-status
```

## Available Commands

```bash
# Development
make init          # Initialize project
make api           # Generate API code with goctl
make swagger       # Generate Swagger documentation
make gen           # Generate API + Swagger docs
make fmt           # Format code
make lint          # Run linter
make test          # Run tests
make build         # Build binary
make run           # Run server
make clean         # Clean build artifacts
make deps          # Install dependencies

# Docker
make docker-build  # Build Docker image
make docker-run    # Run Docker container
make docker-stop   # Stop Docker container
make docker-push   # Push Docker image

# Kubernetes
make k8s-deploy    # Deploy to K8s (default: dev)
make k8s-manifest  # View generated manifest
make k8s-delete    # Delete K8s deployment
make k8s-status    # Check K8s status
```

## API Documentation

### Swagger UI

After starting the service, access Swagger documentation at:
- **JSON**: http://localhost:8888/swagger/swagger.json
- **YAML**: http://localhost:8888/swagger/swagger.yaml

### Module APIs

| Module | Prefix | Description |
|--------|--------|-------------|
| Rule | `/api/v1/rule` | Encoding rule CRUD operations |
| Catalog | `/api/v1/catalog` | Catalog hierarchy management |
| StdFile | `/api/v1/stdfile` | Standard file management |
| Dict | `/api/v1/dict` | Dictionary management |
| DataElement | `/api/v1/dataelement` | Data element operations |
| Task | `/api/v1/task` | Task management |

## Development

### Code Generation

Always use `goctl` to generate API code:

```bash
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Linting

```bash
# Format code
make fmt

# Run linter
make lint
```

## Specification Documents

Each module follows Spec-Driven Development (SDD) workflow:

1. **Context** - Project constitution and guidelines
2. **Specify** - Requirements specification (EARS format)
3. **Design** - Technical design (API/DDL/Model)
4. **Tasks** - Implementation tasks
5. **Implement** - Code implementation

See `specs/` directory for detailed specifications.

## Contributing

1. Follow the [SDD workflow](.specify/workflows/)
2. Create feature branches from `master`
3. Ensure all tests pass before submitting PR
4. Follow Go coding standards

## License

MIT License - see LICENSE file for details

## Contact

- Project Owner: Tian Yuliang
- Repository: https://github.com/tianyuliang/standardization-backend
