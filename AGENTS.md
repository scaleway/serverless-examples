# Agent Instructions: Maintenance and Standardization of Scaleway Serverless Examples

This document defines the strict guidelines that any AI agent (or OpenCode sub-agent) working on the public `scaleway/serverless-examples` repository must follow.

---

## Step 0: Dynamic Runtime Fetching (Mandatory)

Before analyzing or modifying any example:
1. Run the following command in the terminal to query the official Scaleway API if SCW_SECRET_KEY is defined:
   ```bash
   curl -s -H "X-Auth-Token: $SCW_SECRET_KEY" "[https://api.scaleway.com/functions/v1beta1/regions/fr-par/runtimes](https://api.scaleway.com/functions/v1beta1/regions/fr-par/runtimes)"
   ```

2. Extract the list of available runtimes (`status: "available"`) and identify deprecated ones (`status: "deprecated"` or missing).
3. Use this JSON response as the **single source of truth** to validate and update runtimes for each example.

If SCW_SECRET_KEY is not defined or the query fails, consider the following runtimes as latest: node26, python3.14, php8.5, go1.26, rust1.96

---

## Step 1: Deployment Method Standardization

Every repository example MUST offer all supported deployment methods. If a method is missing, create it by faithfully translating the existing configuration (variables, runtime, memory, timeout, triggers).

### 1. Subfolders in `functions/` (3 required methods)
Each subfolder must contain:
- **Terraform:** Working `main.tf` file using the `scaleway/scaleway` provider (`~> 2.0`).
- **Serverless Framework (OSLS):** `serverless.yml` file using the `serverless-scaleway-functions` plugin.
- **Scaleway CLI:** Executable `deploy.sh` script using `scw function ...` commands.

### 2. Subfolders in `containers/` (2 required methods)
Each subfolder must contain:
- **Terraform:** Working `main.tf` file.
- **Scaleway CLI:** Executable `deploy.sh` script using `scw container ...` commands.

### 3. Serverless Framework Requirements (OSLS Usage)
- Prefer the open-source **OSLS** version ([oss-serverless/osls](https://github.com/oss-serverless/osls)) over Serverless v4+.
- In the documentation (`README.md`), specify installation via `npm install -g osls` (or `npx osls`) and use `osls deploy` and `osls remove` for commands.

### 4. CLI Script Requirements (`deploy.sh`)
- Start with `#!/usr/bin/env bash` and `set -euo pipefail`.
- Check for `scw` CLI availability (`command -v scw >/dev/null 2>&1 || { echo "scw required"; exit 1; }`).
- Include creation of namespaces and resources, as well as a commented section at the end showing deletion commands (`scw function namespace delete ...`).

---

## Step 2: Public Repository Standards (Security & Quality)

### 1. Security and Secret Management
- **NO hardcoded secrets:** Never write access keys (`SCW_ACCESS_KEY`), secret keys (`SCW_SECRET_KEY`), tokens, or organization IDs directly in code.
- Read values dynamically from environment variables (`os.Getenv`, `process.env`, `var.scw_secret_key`).
- If the example requires environment variables, include a `.env.example` file containing dummy values.

### 2. Documentation (`README.md`)
Each subfolder must contain a structured `README.md` with:
1. **Description:** What the example does.
2. **Prerequisites:** Required versions for Terraform (`>= 1.5`), OSLS, `scw` CLI, and runtime.
3. **Environment Variables:** List of required variables.
4. **Deployment:** Distinct sections with exact commands for:
   - Terraform (`terraform init && terraform apply`)
   - Serverless Framework (`osls deploy`)
   - Scaleway CLI (`./deploy.sh`)
5. **Cleanup:** Exact commands to destroy created resources.

### 3. Repository Hygiene
- Ensure a proper `.gitignore` file exists.
- Do not commit build artifacts (`.terraform/`, `terraform.tfstate`, `node_modules/`, `.serverless/`, compiled binaries).
- Remove any absolute paths tied to local environments (`/Users/...`, `C:\...`).

---

## Step 3: Execution Mode and Sub-Agents

To process the repository quickly and efficiently without context overflow:

1. **Folder Isolation:** Each sub-agent focuses exclusively on one example subfolder at a time.
2. **Systematic Validation:** If local CLI tools (`terraform`, `scw`, `osls`) are installed, run syntax checks (`terraform fmt -check`, `terraform validate`) before confirming changes.
3. **Reporting:** Each sub-agent summarizes created/modified files and confirms compliance with this document.
