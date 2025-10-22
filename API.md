# API Documentation

## Base URL
```
http://localhost:8080
```

## Common Response Format

### Success Response
```json
{
  "success": true,
  "message": "Success message",
  "data": {},
  "request_id": "unique-request-id"
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error description",
  "request_id": "unique-request-id"
}
```

## Endpoints

### Health Check
Check if the service is running and healthy.

**Endpoint:** `GET /healthcheck`

**Response:**
```json
{
  "success": true,
  "message": "Health check successful",
  "data": {
    "status": "service healthy"
  },
  "request_id": "abc123"
}
```

**Status Codes:**
- `200 OK` - Service is healthy
- `500 Internal Server Error` - Service is unhealthy

**Example:**
```bash
curl http://localhost:8080/healthcheck
```

## Error Handling

All endpoints follow the standard error response format. Common error status codes:

- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Request ID

Every request and response includes a unique `request_id` for tracking and debugging purposes. This ID is also logged in the server logs.
