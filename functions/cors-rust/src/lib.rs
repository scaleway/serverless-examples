use axum::{body::Body, extract::Request, response::Response};
use http::StatusCode;

// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#the_http_response_headers
pub fn with_permissive_cors(r: http::response::Builder) -> http::response::Builder {
    r.header("Access-Control-Allow-Headers", "*")
        .header("Access-Control-Allow-Methods", "*")
        .header("Access-Control-Allow-Origin", "*")
}

pub async fn handler_with_cors(req: Request<Body>) -> Response<Body> {
    println!("{:?}", req);

    with_permissive_cors(Response::builder())
        .status(StatusCode::OK)
        .header("Content-Type", "text/plain")
        .body(Body::from("This is allowing most CORS requests"))
        .unwrap()
}
