# Backend Server

A Go-based backend server with PostgreSQL database integration.

## Setup Instructions

### Prerequisites
- Go 1.23.0 or higher
- PostgreSQL database

### Environment Configuration
1. Copy the example environment file:
   ```bash
   cp env.example .env
   ```

2. Update the `.env` file with your database credentials:
   ```
   DATABASE_URL=postgresql://username:password@localhost:5432/database_name
   PORT=8080
   ```

### Database Setup
1. Create a PostgreSQL database
2. Update the `DATABASE_URL` in your `.env` file with the correct credentials

### Running the Server
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080` (or the port specified in your `.env` file).

### Health Check
Visit `http://localhost:8080/health` to check if the server is running.

### Graceful Shutdown
The server supports graceful shutdown. Use `Ctrl+C` to stop the server safely.
