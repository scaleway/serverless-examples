use axum::{body::{{Body, to_bytes}}, extract::Request, response::Response};
use http::{{method::Method, StatusCode}};

// Reference: https://doc.rust-lang.org/std/iter/trait.Iterator.html#method.product
fn factorial(n: u64) -> u64 {
    (1..=n).product()
}

pub async fn handler(req: Request<Body>) -> Response<Body> {
    if req.method() != Method::POST {
        return Response::builder()
            .status(StatusCode::METHOD_NOT_ALLOWED)
            .header("Content-Type", "text/plain")
            .body(Body::from("Method not allowed"))
            .unwrap();
    }

    // The SQS trigger sends the message content in the body.
    let body = to_bytes(req.into_body(), usize::MAX).await.unwrap();
    let n = match String::from_utf8(body.to_vec()).unwrap().parse::<u64>() {
        Ok(n) => n,
        Err(e) => {
            return Response::builder()
                // Setting the status code to 200 will mark the message as processed.
                .status(StatusCode::OK)
                .header("Content-Type", "text/plain")
                .body(Body::from(format!("Invalid number: {}", e)))
                .unwrap();
        }
    };

    let result = factorial(n);
    println!("rust: factorial of {} is {}", n, result);

    Response::builder()
        // If the status code is not in the 2XX range, the message is considered
        // failed and is retried. In total, there are 3 retries.
        .status(StatusCode::OK)
        .header("Content-Type", "text/plain")
        // Because triggers are asynchronous, the response body is ignored.
        // It's kept here when testing locally.
        .body(Body::from(format!("{}", result)))
        .unwrap()
}
