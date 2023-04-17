"use strict"
import { pathToFileURL } from "url";

const handleCorsPermissive = (event, context, cb) => {
    console.log(event);
    return {
        body: "This is allowing most CORS requests",
        // Permissive configuration
        // See: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#the_http_response_headers
        headers: {
            "Access-Control-Allow-Headers": "*",
            "Access-Control-Allow-Methods": "*",
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Expose-Headers": "*",
            "Content-Type": ["text/plain"],
        }
    }
};

const handleCorsVeryPermissive = (event, context, cb) => {
    console.log(event);
    const headers = event.headers;
    return {
        body: "This is allowing all CORS requests",
        // Very permissive configuration
        // For testing purposes only
        headers: {
            "Access-Control-Allow-Credentials": true,
            "Access-Control-Allow-Headers": headers["Access-Control-Request-Headers"],
            "Access-Control-Allow-Methods": [headers["Access-Control-Request-Method"]],
            "Access-Control-Allow-Origin": headers["Origin"],
            "Access-Control-Expose-Headers": "*",
            "Content-Type": ["text/plain"],
        }
    }
};

export { handleCorsPermissive, handleCorsVeryPermissive };

/* Module was not imported but called directly, so we can test locally.
This will not be executed on Scaleway Functions */
if (import.meta.url === pathToFileURL(process.argv[1]).href) {
  import("@scaleway/serverless-functions").then(scw_fnc_node => {
    scw_fnc_node.serveHandler(handleCorsPermissive, 8080);
  });
}