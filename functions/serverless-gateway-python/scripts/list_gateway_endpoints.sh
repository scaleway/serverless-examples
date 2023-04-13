#!/bin/bash
set -e

curl https://${GATEWAY_HOST}/scw -H 'X-Auth-Token: ${TOKEN}'
