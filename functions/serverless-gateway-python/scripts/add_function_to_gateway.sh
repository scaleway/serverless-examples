#!/bin/bash
set -e

curl -X POST https://${GATEWAY_URL}/scw \
             -H 'X-Auth-Token: ${GATEWAY_TOKEN}' \
             -H 'Content-Type: application/json' \
             -d '{"target":"$1","relative_url":"$2"}'