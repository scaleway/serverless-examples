#!/bin/bash
set -e
SCW_BUCKET=$(terraform output bucket_name)
SCW_NATS_URL=$(terraform output nats_url)

# Nats context creation and selection
nats context save large-messages --server=$SCW_NATS_URL --creds=./large-messages.creds
nats context select large-messages

# Upload file to S3
aws s3 cp $1 s3://$SCW_BUCKET

# Send the name of the file in NATS
nats pub large-messages $(basename $1)
