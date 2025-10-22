# Database Setup Guide

## Prerequisites

- PostgreSQL 12+ installed
- Database user with appropriate permissions

## Configuration

### Environment Variables

Create a `.env` file based on `.env.example`:

```bash
cp .env.example .env
```

Configure the following database variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ewallet_ums
DB_SSL_MODE=disable
```

## Setup Database

### 1. Create Database

```bash
# Login to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE ewallet_ums;

# Create user (optional)
CREATE USER ewallet_user WITH PASSWORD 'your_password';

# Grant privileges
GRANT ALL PRIVILEGES ON DATABASE ewallet_ums TO ewallet_user;

# Exit
\q
```

### 2. Run Migrations

#### Using psql

```bash
# Run UP migration
psql -U postgres -d ewallet_ums -f database/migrations/000001_create_users_table.up.sql

# Run DOWN migration (if needed)
psql -U postgres -d ewallet_ums -f database/migrations/000001_create_users_table.down.sql
```

#### Using migrate tool

Install migrate CLI:

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

Run migrations:

```bash
# Up
migrate -path database/migrations -database "postgresql://postgres:password@localhost:5432/ewallet_ums?sslmode=disable" up

# Down
migrate -path database/migrations -database "postgresql://postgres:password@localhost:5432/ewallet_ums?sslmode=disable" down

# Force version (if stuck)
migrate -path database/migrations -database "postgresql://postgres:password@localhost:5432/ewallet_ums?sslmode=disable" force 1
```

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

### Indexes

- `idx_users_email` - Email lookup
- `idx_users_phone` - Phone lookup
- `idx_users_created_at` - Created at sorting

## Connection Pool Settings

The application uses the following connection pool settings:

- **MaxOpenConns**: 25 - Maximum number of open connections
- **MaxIdleConns**: 25 - Maximum number of idle connections
- **ConnMaxLifetime**: 5 minutes - Maximum lifetime of a connection
- **ConnMaxIdleTime**: 5 minutes - Maximum idle time of a connection

These can be adjusted in `database/postgres.go` based on your load requirements.

## Testing Connection

### Health Check Endpoint

```bash
curl http://localhost:8080/healthcheck
```

Expected response:
```json
{
  "success": true,
  "request_id": "...",
  "message": "Health check successful",
  "data": {
    "status": "healthy"
  }
}
```

### Using psql

```bash
psql -U postgres -d ewallet_ums -c "SELECT version();"
```

## Troubleshooting

### Connection Refused

1. Check if PostgreSQL is running:
   ```bash
   # macOS
   brew services list | grep postgresql
   
   # Linux
   systemctl status postgresql
   ```

2. Check PostgreSQL port:
   ```bash
   netstat -an | grep 5432
   ```

3. Verify pg_hba.conf allows connections from your host

### Authentication Failed

1. Check username and password in `.env`
2. Verify user exists:
   ```sql
   SELECT usename FROM pg_user;
   ```

### Database Does Not Exist

```bash
createdb -U postgres ewallet_ums
```

### Too Many Connections

Adjust PostgreSQL max_connections:

```sql
-- Check current value
SHOW max_connections;

-- Set new value (requires restart)
ALTER SYSTEM SET max_connections = 200;
```

Or adjust application connection pool settings in `database/postgres.go`.

## Backup & Restore

### Backup

```bash
# Backup database
pg_dump -U postgres -d ewallet_ums > backup.sql

# Backup with compression
pg_dump -U postgres -d ewallet_ums | gzip > backup.sql.gz
```

### Restore

```bash
# Restore from backup
psql -U postgres -d ewallet_ums < backup.sql

# Restore from compressed backup
gunzip -c backup.sql.gz | psql -U postgres -d ewallet_ums
```

## Docker Setup

### Using Docker Compose

```bash
# Start PostgreSQL
docker-compose -f docker-compose.dev.yml up -d postgres

# Check logs
docker-compose -f docker-compose.dev.yml logs -f postgres

# Stop
docker-compose -f docker-compose.dev.yml down
```

### Standalone Docker

```bash
# Run PostgreSQL container
docker run --name postgres-ewallet \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=ewallet_ums \
  -p 5432:5432 \
  -d postgres:15-alpine

# Connect to container
docker exec -it postgres-ewallet psql -U postgres -d ewallet_ums
```

## Performance Tips

1. **Enable connection pooling** - Already configured in the application
2. **Use prepared statements** - Reduces parsing overhead
3. **Create appropriate indexes** - Already created for common queries
4. **Monitor slow queries**:
   ```sql
   -- Enable slow query log
   ALTER DATABASE ewallet_ums SET log_min_duration_statement = 1000;
   ```

5. **Regular VACUUM**:
   ```sql
   VACUUM ANALYZE users;
   ```

## Monitoring

### Check Connection Statistics

```sql
SELECT * FROM pg_stat_activity WHERE datname = 'ewallet_ums';
```

### Check Database Size

```sql
SELECT pg_size_pretty(pg_database_size('ewallet_ums'));
```

### Check Table Size

```sql
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

## Migration Best Practices

1. Always test migrations in development first
2. Backup database before running migrations in production
3. Use transactions where possible
4. Keep migrations small and focused
5. Never edit existing migrations - create new ones
6. Document breaking changes

## Resources

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [sqlx Documentation](https://jmoiron.github.io/sqlx/)
- [golang-migrate](https://github.com/golang-migrate/migrate)
