#!/usr/bin/env bash

# set -e

## export DOMAINS=(example.com abc.com)
## export CLOUDFLARE_TOKEN="<token>"
## Usage:   create-records.sh <sub-domain> <dns-content> <dns-type> <proxied>
## Example: create-records.sh usermetrics usermetrics.cognius.net CNAME true

cd "$(dirname "$(dirname "$0")")" || exit 1

main() {
  local name="${1:?}" content="${2:?}" type="${3:-CNAME}" proxied="${4:-true}"

  local domain zone_id
  for domain in "${DOMAINS[@]}"; do
    printf '[INF] Creating record at %-26s' "$domain"
    zone_id="$(cf_get_zone_id "$domain")"
    printf '(%s): ' "$zone_id"

    cf_create_record "$zone_id" "$type" "$name" "$content" "$proxied"
  done
}

cf_verify_token() {
  __call_cf GET "user/tokens/verify" '' '.result.status == "active"' >/dev/null
}

cf_get_zone_id() {
  local domain="$1"
  __call_cf GET "zones?name=$domain" '' '.result[0].id'
}

cf_create_record() {
  local zone_id="$1" type="$2" name="$3" content="$4" proxied="$5"
  __call_cf POST "zones/$zone_id/dns_records" \
    "{\"type\":\"$type\",\"name\":\"$name\",\"content\":\"$content\",\"proxied\":$proxied}" \
    '.result.name'
}

__call_cf() {
  local method="$1" path="$2" data="$3" token="$CLOUDFLARE_TOKEN"
  local query="$4"

  local args=(
    --silent --location
    --request "$method" --url "https://api.cloudflare.com/client/v4/$path"
    --header "Authorization: Bearer $token"
    --header "Content-Type:application/json"
  )

  if test -n "$data"; then
    args+=(--data "$data")
  fi

  local temp_file
  temp_file="$(mktemp)"

  args+=(
    --output "$temp_file"
    --write-out '%{http_code}'
  )

  local status_code
  status_code="$(curl "${args[@]}")"

  if [[ $status_code -lt 200 || $status_code -gt 299 ]] ||
    ! jq --exit-status ".success" "$temp_file" >/dev/null; then
    printf 'ERROR\n'
    printf '  - E%d %s\n' "$status_code" "$(jq --compact-output '.errors[] | [.message]' "$temp_file")" >&2
    printf '  - Response file: %s\n' "$temp_file" >&2
    return 1
  fi

  jq --raw-output --exit-status --monochrome-output "$query" "$temp_file" && rm "$temp_file"
}

# shellcheck disable=SC1091
test -f "$PWD/.env" && source "$PWD/.env"

if ! cf_verify_token; then
  echo "Your token (CLOUDFLARE_TOKEN) is invalid, please use another" >&2
  exit 1
fi

main "$@"
