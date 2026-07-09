const express = require("express");
const { Pool } = require("pg");

const app = express();
const port = process.env.PORT || 8080;

// Retry configuration
const RETRY_INTERVAL_MS = 5000; // Wait between retries
const MAX_RETRIES = 12; // Give up after ~1 minute

// pg reads standard libpq environment variables (PGHOST, PGPORT, etc.)
// which are injected by the Scaleway Container.
const pool = new Pool();

pool.on("error", (err) => {
  console.error("Unexpected error on idle client:", err.message);
});

// SQLSTATE error classes that indicate the connection is gone:
//   08xxx – Connection Exception
//   53xxx – Insufficient Resources
//   57xxx – Operator Intervention
function isConnectionError(err) {
  return err.code && /^(08|53|57)/.test(err.code);
}

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

// Try to establish a connection to the database, retrying until it succeeds
// or MAX_RETRIES is reached. This handles the case where the VPC connection
// is not fully ready when the container starts.
async function waitForDatabase() {
  for (let attempt = 1; attempt <= MAX_RETRIES; attempt++) {
    try {
      await pool.query("SELECT 1");
      console.log("Database connection established.");
      return true;
    } catch (err) {
      console.error(
        `Connection attempt ${attempt}/${MAX_RETRIES} failed: ${err.message}`
      );
      if (attempt < MAX_RETRIES) {
        await sleep(RETRY_INTERVAL_MS);
      }
    }
  }
  return false;
}

app.get("/", (req, res) => {
  res.json({ status: true });
});

app.get("/ping", async (req, res) => {
  try {
    await pool.query("SELECT 1");
    res.json({ status: true });
  } catch (err) {
    console.error("Ping failed:", err.message);
    res.status(500).json({ status: false, error: err.message });
  }
});

async function main() {
  const connected = await waitForDatabase();
  if (!connected) {
    console.error("Could not connect to the database, exiting.");
    process.exit(1);
  }

  app.listen(port, () => {
    console.log(`Server listening on port ${port}`);
  });
}

main();
