"use strict";

import { pathToFileURL } from "url";

/**
 * Compute the factorial of a number
 * @param {number} n 
 * @returns {number} The factorial of n
 */
function factorial(n) {
    let result = 1;
    for (let i = 2; i <= n; i++) {
        result *= i;
    }
    return result;
}

export function handle(event, context, callback) {
    if (event["httpMethod"] !== "POST") {
        return {
            statusCode: 405,
            body: "Method Not Allowed",
            headers: {
                "Content-Type": "text/plain",
            },
        };
    }

    // The SQS trigger sends the message content in the body.
    const n = parseInt(event["body"]);
    console.log(`node: received ${n}`);
    const result = factorial(n);

    console.log(`node: factorial of ${n} is ${result}`);

    return {
        statusCode: 200,
        body: result.toString(),
        headers: {
            "Content-Type": "text/plain",
        },
    };
}

/* Module was not imported but called directly, so we can test locally.
This will not be executed on Scaleway Functions */
if (import.meta.url === pathToFileURL(process.argv[1]).href) {
    import("@scaleway/serverless-functions").then(scw_fnc_node => {
        scw_fnc_node.serveHandler(handle, 8080);
    });
}
