#!/bin/bash

# Example script showing how to use the pgbadger-server API

SERVER_URL="http://localhost:2385"
USERNAME="admin"
PASSWORD="password"  # Change this to match your config

echo "=== pgbadger-server API Example ==="
echo

# 1. Login to get JWT token
echo "1. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST $SERVER_URL/api/v1/login \
  -H "Content-Type: application/json" \
  -d "{\"username\": \"$USERNAME\", \"password\": \"$PASSWORD\"}")

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "Login failed. Check your credentials."
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

echo "Login successful! Token: ${TOKEN:0:20}..."
echo

# 2. Generate a basic report
echo "2. Generating basic report..."
curl -X POST $SERVER_URL/api/v1/report \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "format": "stderr",
    "jobs": 2,
    "verbose": true,
    "top": 10,
    "title": "Sample Report",
    "data_dir": "/var/pv/data/log"
  }' \
  -o report.html

if [ $? -eq 0 ]; then
    echo "Report generated successfully: report.html"
else
    echo "Failed to generate report"
fi
echo

# 3. Generate JSON report with filters
echo "3. Generating JSON report with filters..."
curl -X POST $SERVER_URL/api/v1/report \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "format": "stderr",
    "extension": "json",
    "jobs": 4,
    "top": 20,
    "sample": 5,
    "begin": "2024-01-01T00:00:00Z",
    "end": "2024-12-31T23:59:59Z",
    "dbname": "mydb",
    "exclude_user": ["postgres"],
    "data_dir": "/var/pv/data/log"
  }' \
  -o report.json

if [ $? -eq 0 ]; then
    echo "JSON report generated successfully: report.json"
else
    echo "Failed to generate JSON report"
fi
echo

# 4. Health check
echo "4. Health check..."
curl -s $SERVER_URL/health
echo
echo

echo "=== Example completed ==="