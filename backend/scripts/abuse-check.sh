#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${1:-http://127.0.0.1:8080}"
HITS="${2:-180}"

echo "Abuse check: ${HITS} rapid requests to ${BASE_URL}/api/v1/health"

tmp_file="$(mktemp)"
trap 'rm -f "${tmp_file}"' EXIT

for _ in $(seq 1 "${HITS}"); do
  curl -s -o /dev/null -w "%{http_code}\n" "${BASE_URL}/api/v1/health" >> "${tmp_file}"
done

ok_count="$(grep -c '^200$' "${tmp_file}" || true)"
limited_count="$(grep -c '^429$' "${tmp_file}" || true)"
other_count="$((HITS - ok_count - limited_count))"

echo "200 responses: ${ok_count}"
echo "429 responses: ${limited_count}"
echo "other responses: ${other_count}"

if [ "${limited_count}" -eq 0 ]; then
  echo "WARNING: no rate-limit response detected"
else
  echo "Rate-limit is active"
fi
