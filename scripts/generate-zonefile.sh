#!/usr/bin/env bash

## Generate zonefile from our main 5 domains
## generate-zonefile.sh

export AWS_PROFILE=sts

cd "$(dirname "$(dirname "$0")")" || exit 1

main() {
  local domain
  for domain in "${DOMAINS[@]}"; do
    create_zonefile_from_domain "$domain"
  done
}

get_zone_id() {
  local domain="$1"
  aws route53 list-hosted-zones-by-name \
    --dns-name "$domain" \
    --max-items 1 \
    --no-cli-pager |
    jq --raw-output --monochrome-output '.HostedZones[0].Id' | sed 's|\/hostedzone\/||'
}

create_zonefile_from_domain() {
  local domain="$1" id="" output=""
  local transformer='.ResourceRecordSets[] | "\(.Name) \t\(.TTL) \t\(.Type) \t\(.ResourceRecords[]?.Value)\n"'

  printf 'creating zonefile from "%s" ' "$domain"

  id="$(get_zone_id "$domain")"
  printf '(%s)' "$id"

  output="assets/zonefiles/$domain.zone"
  aws route53 list-resource-record-sets \
    --hosted-zone-id "$id" \
    --no-cli-pager \
    --no-cli-auto-prompt |
    jq \
      --raw-output --monochrome-output \
      --compact-output --join-output \
      "$transformer" >"$output"

  local count
  count="$(wc -l <"$output")"
  printf ': %d\n' "$count"
}

if ! test -d "assets/zonefiles"; then
  mkdir -p "assets/zonefiles"
fi

# shellcheck disable=SC1091
test -f "$PWD/.env" && source "$PWD/.env"
main
