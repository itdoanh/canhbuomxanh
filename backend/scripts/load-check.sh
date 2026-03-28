#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${1:-http://127.0.0.1:8080}"
REQUESTS="${2:-200}"
CONCURRENCY="${3:-20}"

if ! command -v xargs >/dev/null 2>&1; then
  echo "xargs is required"
  exit 1
fi

echo "Load check: ${REQUESTS} requests, concurrency ${CONCURRENCY}, target ${BASE_URL}/api/v1/health"
START_TS="$(date +%s)"

seq 1 "${REQUESTS}" |
  xargs -P "${CONCURRENCY}" -I{} sh -c '
    code="$(curl -s -o /dev/null -w "%{http_code}" "'"${BASE_URL}"'/api/v1/health")"
    if [ "${code}" -ne 200 ]; then
      echo "failed:${code}"
      exit 1
    fi
  '

END_TS="$(date +%s)"
DURATION="$((END_TS - START_TS))"
if [ "${DURATION}" -le 0 ]; then
  DURATION=1
fi

echo "Completed in ${DURATION}s"
echo "Approx RPS: $((REQUESTS / DURATION))"
