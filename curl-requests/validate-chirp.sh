#!/bin/bash
source ./env.sh

URL="$BASE_URL/api/validate_chirp"

echo "POST $URL"
curl -sS -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
  "body": "I had something interesting for breakfast"
}' \
  "$URL" | jq '.'
