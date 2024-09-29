#!/bin/bash
source ./env.sh

URL="$BASE_URL/api/users"

echo "POST $URL"
curl -sS -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
  "email": "user@example.com"
}' \
  "$URL"
