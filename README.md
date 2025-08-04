# Go Library Service

A Go library service, designed to handle books, members, and borrowing operations efficiently.

## Technology Stack

- **Gin**: HTTP web framework used for building RESTful APIs.
- **Gorm**: ORM library for Golang.
- **PostgreSQL**: Storage for data.
- **Redis**: In-memory data structure store used for caching data.
- **Ginkgo**: Testing framework for Go.

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
