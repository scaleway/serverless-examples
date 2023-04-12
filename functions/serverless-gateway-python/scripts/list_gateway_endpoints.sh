#!/bin/bash
set -e

curl https://${GATEWAY_URL}/scw -H 'X-Auth-Token: ${GATEWAY_TOKEN}'