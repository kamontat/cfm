#!/usr/bin/env bash

## First argument: url template (use %s for domain value)
## Example: verify-domain.sh "usermetrics.%s/actuator/health"

export HTTP_CODE_REGEX="200 OK"
export HTTP_SERVER_REGEX="Server: cloudflare"

cd "$(dirname "$(dirname "$0")")" || exit 1

main() {
  local url_template="$1"

  local url domain output
  for domain in "${DOMAINS[@]}"; do
    output=".temp/$domain.curl"
    # shellcheck disable=SC2059
    printf -v url "$url_template" "$domain"
    printf 'Request %-25s: ' "$domain"

    curl -sqI -o "$output" "$url"
    grep -qiE "$HTTP_CODE_REGEX" "$output" && printf '%s' "S" || printf '%s' "E"
    grep -qiE "$HTTP_SERVER_REGEX" "$output" && printf '%s' "S" || printf '%s' "E"
    echo
  done
}

if ! test -d ".temp"; then
  mkdir -p ".temp"
fi

# shellcheck disable=SC1091
test -f "$PWD/.env" && source "$PWD/.env"
main "${1:-%s}"
