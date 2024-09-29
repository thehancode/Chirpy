#!/bin/bash
source ./env.sh

URL="$BASE_URL/api/users"

echo "POST $URL"
curl -sS -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
  "body": "Hello, world!",
  "user_id": "123e4567-e89b-12d3-a456-426614174000"
}' \
  "$URL"
