# Deploying Remote MCP Servers on Scaleway Serverless Containers

A lightweight Go implementation of the **Model Context Protocol (MCP)** deployed as a Scaleway Serverless Container, enabling AI agents (like OpenCode, Claude, and Cursor) to securely invoke tools over HTTP.

The server exposes three tools:

| Tool | Description |
| :--- | :--- |
| `get_server_time` | Returns the current UTC time from the container. |
| `echo` | Echoes back the input message. |
| `get_app_name` | Returns the `SCW_APPLICATION_NAME` environment variable. |

---

## Why Serverless Fits the Model Context Protocol

The Model Context Protocol (MCP) standardizes how LLMs interact with external tools and data sources. Hosting MCP servers on serverless platforms like Scaleway Serverless Containers offers key operational advantages:

| Advantage | Why It Matters for MCP |
| :--- | :--- |
| **Scale-to-Zero Cost** | AI agents invoke tools in intermittent bursts. Serverless ensures you pay only for active execution time, eliminating costs when idle. |
| **Stateless Architecture** | Remote MCP uses HTTP/SSE transports where requests carry their own execution context. Serverless containers scale horizontally without sticky sessions. |
| **Sandbox Security** | Tool execution happens inside isolated, microVM-level container sandboxes, protecting infrastructure from unexpected agent actions. |
| **Zero Maintenance** | Automated OS patching, SSL termination, and horizontal scaling are handled entirely by the cloud provider. |

### Why Scaleway Serverless Containers?

- **Native Container Support:** Run any Docker image (Go, Rust, Python, Node.js) with zero runtime restrictions.
- **Built-in Private Access Control:** Native API Gateway enforcement of `X-Auth-Token` blocks unauthorized traffic before it hits your container — saving compute costs.
- **European Sovereignty & GDPR:** Full data processing compliance in European regions (Paris, Amsterdam, Warsaw).
- **Automated TLS & Custom Domains:** Out-of-the-box HTTPS endpoint provisioning with support for custom domain mapping.

---

## Technical Architecture & Transports

Local MCP servers communicate over standard input/output (`stdio`). Remote serverless MCP deployments require network-based transports over HTTP:

- **Streamable HTTP / SSE:** The MCP server exposes a `/mcp` endpoint accepting POST requests for JSON-RPC 2.0 tool execution and streaming responses back over Server-Sent Events (SSE).
- **Isolated Routing:** This example uses Go's `http.NewServeMux` (Go 1.22+ method routing) rather than the global default router, preventing route hijacking and allowing clean middleware integration.

A `/health` endpoint is also exposed for container health checks.

---

## Prerequisites

- [Terraform](https://www.terraform.io/downloads) `>= 1.5`
- [Scaleway CLI (`scw`)](https://github.com/scaleway/scaleway-cli) (latest)
- [Docker](https://docs.docker.com/get-docker/) (for building and pushing the image)
- [Go](https://go.dev/dl/) `>= 1.22` (only needed for local development/testing)
- A Scaleway account with an [IAM API key](https://www.scaleway.com/en/docs/identity-and-access-management/iam/how-to/create-api-keys/)

---

## Environment Variables

Copy `.env.example` to `.env` and fill in your values:

| Variable | Required | Description |
| :--- | :--- | :--- |
| `SCW_ACCESS_KEY` | Yes | Scaleway IAM access key. |
| `SCW_SECRET_KEY` | Yes | Scaleway IAM secret key. |
| `SCW_APPLICATION_NAME` | No | Value returned by the `get_app_name` tool. Defaults to `scaleway-mcp-go-server` in Terraform, or `[NOT_SET]` at runtime. |
| `SCW_DEFAULT_PROJECT_ID` | No | Target project ID. Uses your default project if omitted. |

```bash
cp .env.example .env
# Edit .env with your real values
```

---

## Deployment

This example offers two deployment methods. **Choose one.**

### Option 1: Terraform

The `main.tf` file provisions the registry namespace, builds and pushes the Docker image, and creates/deployes the container.

```bash
# Export your Scaleway credentials (or set them in a .env file sourced beforehand)
export SCW_ACCESS_KEY="SCWXXXXXXXXXXXXXXXXX"
export SCW_SECRET_KEY="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

terraform init
terraform apply
```

When prompted, type `yes` to confirm. The apply output will print the container URL, the MCP endpoint, and the health endpoint.

To override defaults (e.g. application name or scaling):

```bash
terraform apply -var="application_name=my-app" -var="max_scale=3"
```

### Option 2: Scaleway CLI (`deploy.sh`)

The `deploy.sh` script builds the Docker image, pushes it, and creates/deployes the container using only `scw` commands.

```bash
export SCW_SECRET_KEY="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

chmod +x deploy.sh
./deploy.sh
```

The script prints the container URL, MCP endpoint, and health endpoint on success.

---

## Connecting Clients (e.g. OpenCode)

Add the deployed remote endpoint to your client configuration (such as `opencode.json` or `~/.config/opencode/opencode.json`). Pass your Scaleway IAM secret key via the `X-Auth-Token` header:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "scaleway-tools": {
      "type": "remote",
      "url": "https://<container-name>.<namespace-id>.functions.fnc.fr-par.scw.cloud/mcp",
      "headers": {
        "X-Auth-Token": "YOUR_SCALEWAY_IAM_SECRET_KEY"
      }
    }
  }
}
```

---

## Cleanup

### Terraform

```bash
terraform destroy
```

### Scaleway CLI

The deletion commands are commented out at the bottom of `deploy.sh`. Uncomment and run them, or execute directly:

```bash
scw container container delete region=fr-par <container-id>
scw container namespace delete region=fr-par <namespace-id>
# Optionally, remove the pushed image from the registry:
scw registry namespace delete region=fr-par <registry-namespace-id>
```

---

## Production Best Practices

- **Cold Start Optimization:** Keep container images ultra-minimal (e.g., Go/Rust static binaries inside `distroless` or `alpine`) to ensure cold boots remain under 1 second. For latency-critical agents, set `min_scale = 1`.
- **Database Connection Management:** When connecting tools to databases (like PostgreSQL), use HTTP-based query clients or a connection pooler (e.g., PgBouncer) to prevent connection exhaustion during rapid container scaling.
- **Least-Privilege Tool Execution:** Restrict database user roles to read-only permissions and sanitize all inputs inside tool implementations to prevent SQL injection or accidental data manipulation.
