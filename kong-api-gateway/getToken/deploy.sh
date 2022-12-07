#!/bin/bash
test -f token && echo "Token already exists, re-using." && exit 0
echo "Generating token in namespace $NAMESPACE_ID in region $REGION_NAME"
curl -sSL -X POST -H "X-Auth-Token: $SCW_SECRET_KEY" -H 'Content-Type: application/json' -d "{\"namespace_id\": \"$NAMESPACE_ID\"}"  https://api.scaleway.com/functions/v1beta1/regions/$REGION_NAME/tokens | jq .token | sed 's/"//g' > token
