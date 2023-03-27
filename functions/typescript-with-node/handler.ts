export {handle};

function handle (event: Record<string, unknown>, context: Record<string, unknown>, cb: unknown) {
    return {
        body: "Hello world!",
        headers: { "Content-Type": ["application/json"] },
        statusCode: 200,
    };
};