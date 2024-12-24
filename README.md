## TrinityKnights Backend

TrinityKnights Backend is a service that provides backend functionalities for the TrinityKnights application.

[![Go](https://img.shields.io/github/go-mod/go-version/savioruz/TrinityKnights.Backend)](https://golang.org/)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/savioruz/TrinityKnights.Backend)
[![Go Report Card](https://goreportcard.com/badge/github.com/savioruz/TrinityKnights.Backend)](https://goreportcard.com/report/github.com/savioruz/TrinityKnights.Backend)
![License](https://img.shields.io/github/license/savioruz/TrinityKnights.Backend)

## Table of Contents
- [Features](#features)
- [Deployment](#deployment)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
  - [Running the API](#running-the-api)
    - [Docker](#docker)
    - [Make](#make)
  - [API Documentation](#api-documentation)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)
- [Reference](#reference)
- [Acknowledgements](#acknowledgements)

## Features

- **Authentication & Authorization**
  - JWT-based authentication
  - Role-based access control (Admin/Public)
  - Token refresh mechanism

- **Payment Processing**
  - Xendit payment gateway integration
  - Payment status tracking
  - Invoice generation
  - Payment search and filtering

- **API Interfaces**
  - RESTful HTTP API
  - GraphQL API with playground
  - Swagger/OpenAPI documentation

- **Data Management**
  - PostgreSQL database integration
  - Redis caching
  - GORM ORM implementation

- **Development Features**
  - Docker and Docker Compose support
  - Comprehensive test coverage
  - Mock generation for testing
  - Code quality checks (critic, security)
  - Hot reload in development

## Deployment

- ### Koyeb
[![Deploy to Koyeb](https://www.koyeb.com/static/images/deploy/button.svg)](https://app.koyeb.com/services/deploy?type=git&builder=dockerfile&repository=github.com/savioruz/TrinityKnights.Backend&branch=main&ports=3000;http;/&name=TrinityKnights-Backend&env[STAGE_STATUS]=prod&env[APP_NAME]=TrinityKnights-Backend&env[APP_HOST]=0.0.0.0&env[APP_PORT]=3000&env[DB_HOST]=YOUR_DB_HOST&env[DB_PORT]=5432&env[DB_USER]=YOUR_DB_USER&env[DB_PASSWORD]=YOUR_DB_PASSWORD&env[DB_NAME]=YOUR_DB_NAME&env[DB_SSL_MODE]=require&env[DB_TIMEZONE]=YOUR_DB_TIMEZONE&env[REDIS_HOST]=YOUR_REDIS_HOST&env[REDIS_PORT]=6379&env[REDIS_PASSWORD]=&env[REDIS_DB]=0)

- ### Railway
[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/jT1IvF?referralCode=XVMtOY)

- ### Vercel
[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fsavioruz%2FTrinityKnights.Backend&env=STAGE_STATUS,APP_NAME,APP_HOST,APP_PORT,DB_HOST,DB_PORT,DB_USER,DB_PASSWORD,DB_NAME,DB_SSL_MODE,DB_TIMEZONE,REDIS_HOST,REDIS_PORT,REDIS_PASSWORD,REDIS_DB&envDescription=env&project-name=TrinityKnights-Backend&repository-name=TrinityKnights-Backend)

- ### Render
[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/savioruz/TrinityKnights.Backend)

## Requirements

- Go 1.23+
- Docker
- PostgreSQL
- Redis
- Make

## Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/savioruz/TrinityKnights.Backend.git
    cd TrinityKnights.Backend
    ```

2. **Environment Variables:**

   Create a `.env` file in the root directory and add the following:

    ```bash
    cp .env.example .env
    ```

3. **Edit the `.env` file with your configuration:**

    ```dotenv
    APP_PORT=3000
    APP_LOG_LEVEL=6
    APP_ENV=development

    DB_HOST=
    DB_PORT=5432
    DB_USER=
    DB_PASSWORD=
    DB_NAME=
    DB_SSL_MODE=require
    DB_TIMEZONE=

    REDIS_HOST=
    REDIS_PORT=
    REDIS_PASSWORD=
    REDIS_DB=0

    SMTP_HOST=
    SMTP_PORT=
    SMTP_USERNAME=
    SMTP_PASSWORD=

    JWT_SECRET=
    JWT_ACCESS_EXPIRY=1h
    JWT_REFRESH_EXPIRY=168h

    XENDIT_API_KEY=
    XENDIT_CALLBACK_TOKEN=
    ```

## Usage

### Development Commands

```bash
# Generate Swagger documentation
make swag

# Generate mocks for testing
make mockgen

# Run code critics (static analysis)
make critic

# Run security checks
make security

# Run tests with coverage
make test
```

### Docker Commands

#### Single Container
```bash
# Build Docker image
make docker.build

# Run container
make docker.run

# Stop and remove container
make docker.stop
```

#### Docker Compose
```bash
# Build services
make dc.build

# Start services
make dc.up

# Stop services
make dc.down
```

### Testing

Run the full test suite with coverage:
```bash
make test
```
This will:
1. Clean build artifacts
2. Run code critics
3. Perform security checks
4. Execute tests with coverage
5. Display coverage report

### Code Generation

Generate mock implementations for testing:
```bash
make mockgen
```

Generate Swagger documentation:
```bash
make swag
```

### Running the API

You can run the API using Docker or directly with Make.

### Docker

1. **Run redis:**

    ```bash
    make docker.redis
    ```

2. **Run the application:**

    ```bash
    make docker.run
    ```

For production, you need to secure redis on Makefile with a password.

### Make

1. **Run the application:**

    ```bash
    make run
    ```

You need to have Redis running on your machine.

### ERD Schema

![ERD Schema](/assets/erd.png)

### API Documentation

Swagger documentation is available at: http://localhost:3000/swagger.

### GraphQL Playground

GraphQL playground is available at: http://localhost:3000/playground.

## Project Structure

```
.
├── cmd/
│   └── app/main.go
├── config/
├── docs/
├── internal/
│   ├── builder/
│   ├── delivery/
│   │   ├── graph/
│   │   │   ├── handler/
│   │   │   ├── model/
│   │   │   ├── resolver/
│   │   │   └── schema.graphqls
│   │   └── http/
│   │   │   ├── handler/
│   │   │   ├── middleware/
│   │   │   └── route/
│   ├── domain/
│   │   ├── entity/
│   │   └── model/
│   │       └── converter/
│   ├── repository/
│   └── service/
├── pkg/
│   ├── cache/
│   ├── errors/
│   ├── helper/
│   └── route/
├── test/mock
├── .env
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── tools.go
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── vercel.json
└── LICENSE
    

```

## Contributing

Feel free to open issues or submit pull requests with improvements.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Reference

- [Go Programming Language](https://golang.org/)
- [GORM Documentation](https://gorm.io/docs/)
- [GraphQL](https://graphql.org/)
- [Xendit API Documentation](https://developers.xendit.co/)
- [Redis Documentation](https://redis.io/documentation)
- [Docker Documentation](https://docs.docker.com/)
- [Swagger/OpenAPI](https://swagger.io/specification/)

## Acknowledgements

- [Go-Swagger](https://github.com/go-swagger/go-swagger)
- [GQLGen](https://gqlgen.com/)
- [GoMock](https://github.com/golang/mock)
- [Logrus](https://github.com/sirupsen/logrus)
- [Go-Playground Validator](https://github.com/go-playground/validator)
- [Xendit Go Library](https://github.com/xendit/xendit-go)
