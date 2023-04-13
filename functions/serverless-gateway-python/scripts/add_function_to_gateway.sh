#!/bin/bash
set -e

curl -X POST https://${GATEWAY_HOST}/scw \
             -H 'X-Auth-Token: ${TOKEN}' \
             -H 'Content-Type: application/json' \
             -d '{"target":"$1","relative_url":"$2"}'
