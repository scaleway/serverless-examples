"""Minimal FastAPI streaming example for debugging purposes."""

import asyncio
import json
import time
import logging
import os

from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import StreamingResponse, JSONResponse
from fastapi.staticfiles import StaticFiles

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(title="FastAPI Streaming Example")

# Mount static files
_static_dir = os.path.join(os.path.dirname(__file__), "static")
app.mount("/static", StaticFiles(directory=_static_dir), name="static")

# Allow CORS for browser-based clients
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=False,
    allow_methods=["GET", "POST", "OPTIONS"],
    allow_headers=["*"],
)

# --- SSE helpers -----------------------------------------------------------

def sse_event(data: dict) -> bytes:
    """Format a dict as a Server-Sent Event."""
    return f"data: {json.dumps(data)}\n\n".encode("utf-8")


# --- Streaming generator ----------------------------------------------------

async def event_stream(chunks: int = 10, delay: float = 1.0):
    """Yield ``chunks`` SSE events with ``delay`` seconds between each."""
    for i in range(chunks):
        payload = {
            "id": i,
            "content": f"Chunk {i} of {chunks}",
            "timestamp": time.time(),
        }
        logger.info("Sending chunk %d/%d", i + 1, chunks)
        yield sse_event(payload)
        await asyncio.sleep(delay)

    # Final event
    yield sse_event({"done": True, "total_chunks": chunks})


# --- Routes -----------------------------------------------------------------

@app.get("/")
async def root():
    """Serve the test HTML page."""
    from fastapi.responses import FileResponse
    return FileResponse(os.path.join(_static_dir, "index.html"))


@app.get("/health")
async def health():
    return {"status": "healthy"}


@app.options("/public/completions")
async def completions_options():
    """Handle CORS preflight for the streaming endpoint."""
    return JSONResponse(content={}, headers={
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "POST, OPTIONS",
        "Access-Control-Allow-Headers": "*",
    })


@app.post("/public/completions")
async def completions(request: Request):
    """Stream a configurable number of SSE events.

    Query params:
        chunks: number of SSE events to emit (default 10)
        delay:  seconds between each event (default 1.0)
    """
    params = request.query_params
    chunks = int(params.get("chunks", 10))
    delay = float(params.get("delay", 1.0))

    logger.info(
        "Starting stream: chunks=%d delay=%.2f client=%s",
        chunks, delay, request.client.host if request.client else "unknown",
    )

    return StreamingResponse(
        event_stream(chunks=chunks, delay=delay),
        media_type="text/event-stream",
        headers={
            "Cache-Control": "no-cache",
            "Connection": "keep-alive",
            "X-Accel-Buffering": "no",  # disable nginx buffering
        },
    )


# --- Long stream for stress testing ----------------------------------------

@app.post("/public/completions/long")
async def completions_long(request: Request):
    """Stream events for a long duration to test timeouts.

    Query params:
        duration: total seconds to stream (default 300)
        interval: seconds between events (default 5)
    """
    params = request.query_params
    duration = float(params.get("duration", 300))
    interval = float(params.get("interval", 5))

    async def long_stream():
        start = time.time()
        i = 0
        while time.time() - start < duration:
            i += 1
            elapsed = time.time() - start
            yield sse_event({
                "id": i,
                "elapsed": round(elapsed, 2),
                "remaining": round(max(0, duration - elapsed), 2),
            })
            await asyncio.sleep(interval)

        yield sse_event({"done": True, "total_events": i})

    logger.info(
        "Starting long stream: duration=%.0f interval=%.2f client=%s",
        duration, interval, request.client.host if request.client else "unknown",
    )

    return StreamingResponse(
        long_stream(),
        media_type="text/event-stream",
        headers={
            "Cache-Control": "no-cache",
            "Connection": "keep-alive",
            "X-Accel-Buffering": "no",
        },
    )


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
