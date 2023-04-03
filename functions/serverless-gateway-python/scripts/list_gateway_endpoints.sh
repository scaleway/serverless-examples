#!/bin/bash
set -e

curl http://${GATEWAY_URL}/scw -H 'X-Auth-Token: ${GATEWAY_TOKEN}'