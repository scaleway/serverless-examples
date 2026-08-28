#!/usr/bin/env bash
set -euo pipefail

# =============================================================================
# Scaleway CLI deployment script for the mcp-go-server container example.
# Builds the Docker image, pushes it to the Scaleway Container Registry, then
# creates/deployes a Serverless Container exposing the MCP server.
# =============================================================================

# --- Configuration -----------------------------------------------------------
NAMESPACE_NAME="mcp-go-server-ns"
CONTAINER_NAME="mcp-go-server"
IMAGE_NAME="mcp-server"
IMAGE_TAG="latest"
PORT=8080
MIN_SCALE=0
MAX_SCALE=5
PRIVACY="private"
REGION="fr-par"
REGISTRY="rg.${REGION}.scw.cloud"

# --- Pre-flight checks -------------------------------------------------------
if ! command -v scw >/dev/null 2>&1; then
  echo "Error: Scaleway CLI ('scw') is not installed or not in PATH."
  echo "Install it from: https://github.com/scaleway/scaleway-cli"
  exit 1
fi

if ! command -v docker >/dev/null 2>&1; then
  echo "Error: 'docker' is required but was not found in PATH."
  exit 1
fi

if ! command -v jq >/dev/null 2>&1; then
  echo "Error: 'jq' is required but was not found in PATH."
  exit 1
fi

if [[ -z "${SCW_SECRET_KEY:-}" ]]; then
  echo "Error: SCW_SECRET_KEY environment variable is not set."
  echo "Export your Scaleway IAM secret key before running this script."
  exit 1
fi

# Use a per-run unique registry namespace name to avoid collisions.
# Scaleway registry namespace names must be unique per region.
REGISTRY_NAMESPACE="mcp-go-server-$(date +%s | tail -c 5)"

# --- Docker registry login ---------------------------------------------------
echo "Logging in to Scaleway Container Registry..."
echo "${SCW_SECRET_KEY}" | docker login "${REGISTRY}" -u nologin --password-stdin

# --- Registry namespace ------------------------------------------------------
# The registry namespace must exist before the image can be pushed.
echo "Ensuring the registry namespace '${REGISTRY_NAMESPACE}' exists..."
REGISTRY_NAMESPACE_ID="$(scw registry namespace list \
  region="${REGION}" \
  name="${REGISTRY_NAMESPACE}" \
  -o json | jq -r '.[0].id // empty')"

if [[ -z "${REGISTRY_NAMESPACE_ID}" ]]; then
  echo "Creating registry namespace '${REGISTRY_NAMESPACE}'..."
  REGISTRY_NAMESPACE_ID="$(scw registry namespace create \
    region="${REGION}" \
    name="${REGISTRY_NAMESPACE}" \
    -o json | jq -r '.id')"
fi

# --- Build and push the image ------------------------------------------------
FULL_IMAGE="${REGISTRY}/${REGISTRY_NAMESPACE}/${IMAGE_NAME}:${IMAGE_TAG}"

echo "Building Docker image: ${FULL_IMAGE}"
docker build -t "${FULL_IMAGE}" .

echo "Pushing Docker image..."
docker push "${FULL_IMAGE}"

# --- Container namespace -----------------------------------------------------
NAMESPACE_ID="$(scw container namespace list \
  region="${REGION}" \
  name="${NAMESPACE_NAME}" \
  -o json | jq -r '.[0].id // empty')"

if [[ -z "${NAMESPACE_ID}" ]]; then
  echo "Creating container namespace '${NAMESPACE_NAME}'..."
  NAMESPACE_ID="$(scw container namespace create \
    region="${REGION}" \
    name="${NAMESPACE_NAME}" \
    --wait \
    -o json | jq -r '.id')"
else
  echo "Reusing existing namespace '${NAMESPACE_NAME}' (${NAMESPACE_ID})."
fi

# --- Container ---------------------------------------------------------------
echo "Creating container '${CONTAINER_NAME}'..."
CONTAINER_ID="$(scw container container create \
  region="${REGION}" \
  namespace-id="${NAMESPACE_ID}" \
  name="${CONTAINER_NAME}" \
  image="${FULL_IMAGE}" \
  port="${PORT}" \
  min-scale="${MIN_SCALE}" \
  max-scale="${MAX_SCALE}" \
  privacy="${PRIVACY}" \
  --wait \
  -o json | jq -r '.id')"

# --- Output ------------------------------------------------------------------
DOMAIN="$(scw container container get \
  region="${REGION}" \
  "${CONTAINER_ID}" \
  -o json | jq -r '.public_endpoint')"

echo ""
echo "Container URL:    ${DOMAIN}"
echo "MCP endpoint:     ${DOMAIN}/mcp"
echo "Health endpoint:  ${DOMAIN}/health"

# =============================================================================
# Cleanup / Remove
# -----------------------------------------------------------------------------
# To delete the container and namespace created above, uncomment and run:
#
#   scw container container delete region="${REGION}" "${CONTAINER_ID}"
#   scw container namespace delete region="${REGION}" "${NAMESPACE_ID}"
#
# To also remove the pushed Docker image from the registry:
#
#   scw registry namespace delete region="${REGION}" "${REGISTRY_NAMESPACE_ID}"
# =============================================================================
