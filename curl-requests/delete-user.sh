#!/bin/bash
source ./env.sh

URL="$BASE_URL/admin/reset"

echo "POST $URL"
curl -sS -X POST \
  "$URL"
