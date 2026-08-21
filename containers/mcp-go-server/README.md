# Deploying Remote MCP Servers on Scaleway Serverless Containers

A lightweight guide on how and why to run Model Context Protocol (MCP) servers on Scaleway Serverless Containers, enabling AI agents (like OpenCode, Claude, and Cursor) to securely invoke tools over HTTP.

---

## Why Serverless Fits the Model Context Protocol

The Model Context Protocol (MCP) standardizes how LLMs interact with external tools and data sources. Hosting MCP servers on serverless platforms like Scaleway Serverless Containers or GCP Cloud Run offers key operational advantages:

| Advantage | Why It Matters for MCP |
| :--- | :--- |
| **Scale-to-Zero Cost** | AI agents invoke tools in intermittent bursts. Serverless ensures you pay only for active execution time, eliminating costs when idle. |
| **Stateless Architecture** | Remote MCP uses HTTP/SSE transports where requests carry their own execution context. Serverless containers scale horizontally without sticky sessions. |
| **Sandbox Security** | Tool execution happens inside isolated, microVM-level container sandboxes, protecting infrastructure from unexpected agent actions. |
| **Zero Maintenance** | Automated OS patching, SSL termination, and horizontal scaling are handled entirely by the cloud provider. |

---

## Why Scaleway Serverless Containers?

* **Native Container Support:** Run any Docker image (Go, Rust, Python, Node.js) with zero runtime restrictions.
* **Built-in Private Access Control:** Native API Gateway enforcement of `X-Auth-Token` blocks unauthorized traffic before it hits your container—saving compute costs.
* **European Sovereignty & GDPR:** Full data processing compliance in European regions (Paris, Amsterdam, Warsaw).
* **Automated TLS & Custom Domains:** Out-of-the-box HTTPS endpoint provisioning with support for custom domain mapping.

---

## Technical Architecture & Transports

Local MCP servers communicate over standard input/output (`stdio`). Remote serverless MCP deployments require network-based transports over HTTP:

* **Streamable HTTP / SSE:** The MCP server exposes a `/mcp` endpoint accepting POST requests for JSON-RPC 2.0 tool execution and streaming responses back over Server-Sent Events (SSE).
* **Isolated Routing:** Production implementations use isolated routers (like Go's `http.NewServeMux`) rather than default global routers to prevent route hijacking and allow clean middleware integration.


---

## Deployment Steps

### 1. Build and Push Container Image

Containerize your MCP server binary and push it to the Scaleway Container Registry:

```bash
# Login to Scaleway Container Registry
docker login rg.fr-par.scw.cloud -u nologin -p $SCW_SECRET_KEY

# Build and push your lightweight container image
docker build -t rg.fr-par.scw.cloud/<namespace>/<image_name>:v1 .
docker push rg.fr-par.scw.cloud/<namespace>/<image_name>:v1

```

### 2. Deploy to Scaleway Serverless

Deploy the container using the Scaleway CLI or console. For production, set privacy to `private` to enforce IAM token authentication:

```bash
# Create the container deployment
scw container container create \
  name=remote-mcp-server \
  namespace-id=<your-namespace-id> \
  image=rg.fr-par.scw.cloud/<namespace>/<image_name>:v1 \
  port=8080 \
  min-scale=0 \
  max-scale=5 \
  privacy=private

# Deploy the instance
scw container container deploy <container-id>

```

---

## Connecting Clients (e.g., OpenCode)

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

## Production Best Practices

* **Cold Start Optimization:** Keep container images ultra-minimal (e.g., Go/Rust static binaries inside `distroless` or `alpine`) to ensure cold boots remain under 1 second. For latency-critical agents, set `--min-scale=1`.
* **Database Connection Management:** When connecting tools to databases (like PostgreSQL), use HTTP-based query clients or a connection pooler (e.g., PgBouncer) to prevent connection exhaustion during rapid container scaling.
* **Least-Privilege Tool Execution:** Restrict database user roles to read-only permissions and sanitize all inputs inside tool implementations to prevent SQL injection or accidental data manipulation.
