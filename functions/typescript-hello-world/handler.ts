export {handle};

function handle (event, context, cb) {
    return {
        body: "Hello world!",
        headers: { "Content-Type": ["application/json"] },
        statusCode: 200,
    };
};