#!/bin/bash

set -e

echo "Connecting to ${MONGO_HOSTNAME}"

mongosh mongodb+srv://${MONGO_HOSTNAME} --apiVersion 1 --username ${MONGO_USERNAME} --password ${MONGO_PASSWORD} --eval "db.version()"
