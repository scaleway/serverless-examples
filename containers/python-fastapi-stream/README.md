# FastAPI Streaming Example

A minimal FastAPI container with Server-Sent Events (SSE) streaming responses, designed for debugging network issues on Scaleway Serverless Containers.

## Endpoints

| Method    | Path                       | Description                              |
| --------- | -------------------------- | ---------------------------------------- |
| `GET`     | `/`                        | Browser test page (HTML)                 |
| `GET`     | `/health`                  | Health check                             |
| `POST`    | `/public/completions`      | Stream SSE events (configurable)         |
| `POST`    | `/public/completions/long` | Long-duration stream for timeout testing |
| `OPTIONS` | `/public/completions`      | CORS preflight                           |

### `/public/completions` query params

| Param    | Default | Description                  |
| -------- | ------- | ---------------------------- |
| `chunks` | `10`    | Number of SSE events to emit |
| `delay`  | `1.0`   | Seconds between each event   |

### `/public/completions/long` query params

| Param      | Default | Description             |
| ---------- | ------- | ----------------------- |
| `duration` | `300`   | Total seconds to stream |
| `interval` | `5`     | Seconds between events  |

## Local development (Nix)

```bash
nix-shell
pip install -r requirements.txt
uvicorn app.main:app --host 0.0.0.0 --port 8080
```

## Build & deploy (Scaleway)

```bash
scw container deploy
```

## Testing

### Browser

Open `http://localhost:8080/` (or your deployed container URL) in a browser. The test page lets you:

- Set the number of chunks and delay between events
- Start, stop, and clear streams
- See SSE events and network errors in real time

### CLI

```bash
# Quick test
curl -N -X POST http://localhost:8080/public/completions?chunks=5&delay=1

# Long stream (test timeouts)
curl -N -X POST http://localhost:8080/public/completions/long?duration=600&interval=5
```

## Streaming headers

The response includes:

- `Content-Type: text/event-stream`
- `Cache-Control: no-cache`
- `Connection: keep-alive`
- `X-Accel-Buffering: no` (disables nginx proxy buffering)
