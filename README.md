# pgbadger-server

A production-ready HTTP server that wraps pgbadger with secure authentication and configurable options.

## Features

- **JWT Authentication**: Secure username/password authentication with JWT tokens
- **Whitelist Security**: Only allows predefined pgbadger options to prevent command injection
- **Parameter Validation**: Comprehensive validation of all input parameters
- **Timeout Protection**: 5-minute timeout on pgbadger execution
- **Structured Logging**: Detailed audit logs for security and debugging

## Quick Start

### 1. Configure Authentication

Edit `config/config.yaml`:

```yaml
auth:
  jwt_secret: "your-super-secret-jwt-key-change-this-in-production"
  jwt_expiry: "24h"
  username: "admin"
  password_hash: "$2a$10$example.hash.here"  # Use bcrypt hash of your password
```

### 2. Generate Password Hash

```bash
# Install bcrypt tool or use online generator
go run -c 'package main; import ("fmt"; "golang.org/x/crypto/bcrypt"); func main() { hash, _ := bcrypt.GenerateFromPassword([]byte("yourpassword"), 10); fmt.Println(string(hash)) }'
```

### 3. Run the Server

```bash
go run main.go server
```

## API Usage

### Authentication

First, obtain a JWT token:

```bash
curl -X POST http://localhost:2385/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "yourpassword"}'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Generate Report

Use the token to generate reports:

```bash
curl -X POST http://localhost:2385/api/v1/report \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "format": "stderr",
    "jobs": 4,
    "verbose": true,
    "top": 20,
    "title": "My PostgreSQL Report",
    "data_dir": "/var/log/postgresql"
  }'
```

## Supported pgbadger Options

The API accepts a JSON object with the following whitelisted options:

### Basic Options
- `format`: Log format (`"syslog"`, `"syslog2"`, `"stderr"`, `"jsonlog"`, `"csv"`, `"pgbouncer"`, `"logplex"`, `"rds"`, `"redshift"`)
- `outfile`: Output filename
- `outdir`: Output directory
- `title`: Report title
- `jobs`: Number of parallel jobs (1-32)
- `jobs_parallel`: Number of files to parse in parallel (1-32)
- `verbose`: Enable verbose output
- `quiet`: Suppress output

### Filtering Options
- `dbname`: Filter by database name
- `dbuser`: Filter by database user
- `appname`: Filter by application name
- `client_host`: Filter by client host

### Time Options
- `begin`: Start datetime (RFC3339 or `2006-01-02 15:04:05`)
- `end`: End datetime (RFC3339 or `2006-01-02 15:04:05`)

### Report Options
- `top`: Number of queries to display (1-1000)
- `sample`: Number of query samples (1-100)
- `maxlength`: Maximum query length (100-1000000)

### Graph Options
- `average`: Minutes for average graphs (1-1440)
- `histo_average`: Minutes for histogram graphs (1-10080)
- `nograph`: Disable graphs

### Output Options
- `extension`: Output format (`"text"`, `"html"`, `"bin"`, `"json"`)
- `prettify`: Enable SQL prettification (default: true)
- `query_numbering`: Add query numbering

### Special Modes
- `select_only`: Report only SELECT queries
- `watch_mode`: Report only errors
- `incremental`: Use incremental mode
- `explode`: Generate separate reports per database

### Exclude Options (Arrays)
- `exclude_user`: Users to exclude
- `exclude_appname`: Application names to exclude
- `exclude_client`: Client IPs to exclude
- `exclude_db`: Databases to exclude

### Include Options (Arrays)
- `include_query`: Regex patterns for queries to include
- `include_pid`: Process IDs to include
- `include_session`: Session IDs to include

### Advanced Options
- `prefix`: Custom log_line_prefix
- `ident`: Syslog ident
- `timezone`: Timezone offset
- `log_timezone`: Log timezone offset

### Data Directory
- `data_dir`: Path to log files (defaults to `/var/pv/data/log`)

## Security Features

- **JWT Authentication**: Bearer token required for all API endpoints
- **Parameter Whitelisting**: Only predefined options are accepted
- **Input Validation**: Strict type checking and range validation
- **Command Injection Prevention**: Uses `exec.Command()` with argument arrays
- **Timeout Protection**: 5-minute execution timeout
- **Audit Logging**: All requests and command executions are logged

## Configuration

The server supports configuration via:

1. **YAML Config File** (`config/config.yaml`)
2. **Environment Variables** (prefixed with section, e.g., `AUTH_JWT_SECRET`)
3. **Command Line Flags** (when implemented)

## Production Deployment

1. **Change JWT Secret**: Use a strong, random secret
2. **Use HTTPS**: Deploy behind a reverse proxy with SSL
3. **Secure Passwords**: Use strong passwords with bcrypt hashing
4. **Log Rotation**: Configure log rotation for audit trails
5. **Resource Limits**: Set appropriate memory and CPU limits
6. **Monitoring**: Implement health checks and metrics

## Health Check

```bash
curl http://localhost:2385/health
```

Returns `{"status": "ok"}` when healthy.