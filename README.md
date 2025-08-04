# Go Library Service

A Go-based library service for managing books, members, and borrow.

## Prerequisites

- Makefile
- Docker and Docker Compose

## Getting Started

1. **Clone the repository**

   ```bash
   git clone https://github.com/Kittiphop/go-library-service.git
   cd go-library-service
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```
   Edit the `.env` file to configure your environment settings.

## Development

### Start Development Environment

To start the local development environment with Docker:

```bash
make dev.up
```

This will:

- Start all required services using Docker Compose
- Swagger UI: http://localhost:8080/swagger/index.html
- Staff account for testing

```
username: staff
password: staff
```

### Stop Development Environment

To stop the local development environment:

```bash
make dev.down
```

This will stop and remove all containers and volumes.

## Testing

To run the test suite:

```bash
make test
```

This will run all unit tests using ginkgo.
