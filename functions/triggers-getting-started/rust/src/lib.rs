use hyper::{Body, Request, Response, StatusCode};

// Reference: https://doc.rust-lang.org/std/iter/trait.Iterator.html#method.product
fn factorial(n: u64) -> u64 {
    (1..=n).product()
}

pub async fn handler(req: Request<Body>) -> Response<Body> {
    if req.method() != hyper::Method::POST {
        return Response::builder()
            .status(StatusCode::METHOD_NOT_ALLOWED)
            .header("Content-Type", "text/plain")
            .body(Body::from("Method not allowed"))
            .unwrap();
    }

    // The SQS trigger sends the message content in the body.
    let body = hyper::body::to_bytes(req.into_body()).await.unwrap();
    let n = match String::from_utf8(body.to_vec()).unwrap().parse::<u64>() {
        Ok(n) => n,
        Err(e) => {
            return Response::builder()
                .status(StatusCode::BAD_REQUEST)
                .header("Content-Type", "text/plain")
                .body(Body::from(format!("Invalid number: {}", e)))
                .unwrap();
        }
    };

    let result = factorial(n);
    println!("rust: factorial of {} is {}", n, result);

    Response::builder()
        .status(StatusCode::OK)
        .header("Content-Type", "text/plain")
        .body(Body::from(format!("{}", result)))
        .unwrap()
}
