use axum::{
    body::Body, extract::FromRequest, extract::Multipart, extract::Request, response::Response,
};
use http::StatusCode;

// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#the_http_response_headers
pub fn with_permissive_cors(r: http::response::Builder) -> http::response::Builder {
    r.header("Access-Control-Allow-Headers", "*")
        .header("Access-Control-Allow-Methods", "*")
        .header("Access-Control-Allow-Origin", "*")
}

pub async fn handler_with_cors(req: Request<Body>) -> Response<Body> {
    println!("{:?}", req);

    // Check if this is an OPTIONS request
    if req.method() == http::Method::OPTIONS {
        return with_permissive_cors(Response::builder())
            .status(StatusCode::OK)
            .header("Content-Type", "text/plain")
            .body(Body::from("This is allowing most CORS requests"))
            .unwrap();
    }

    // For POST requests, extract multipart data
    let (parts, body) = req.into_parts();
    let body_stream = axum::body::Body::new(body);

    // Reconstruct the request for multipart extraction
    let request = Request::from_parts(parts, body_stream);

    // Use Axum's built-in multipart extractor
    match Multipart::from_request(request, &()).await {
        Ok(multipart) => process_multipart(multipart).await,
        Err(err) => with_permissive_cors(Response::builder())
            .status(StatusCode::BAD_REQUEST)
            .body(Body::from(format!(
                "Failed to extract multipart data: {}",
                err
            )))
            .unwrap(),
    }
}

// Process the multipart data, using Axum's multipart exmaple
// See: https://github.com/tokio-rs/axum/blob/main/examples/multipart-form/src/main.rs
async fn process_multipart(mut multipart: Multipart) -> Response<Body> {
    while let Some(field) = multipart.next_field().await.unwrap() {
        let name = field.name().unwrap().to_string();
        let file_name = field.file_name().unwrap().to_string();
        let content_type = field.content_type().unwrap().to_string();
        let data = field.bytes().await.unwrap();

        println!(
            "Length of `{name}` (`{file_name}`: `{content_type}`) is {} bytes",
            data.len()
        );
    }

    with_permissive_cors(Response::builder())
        .status(StatusCode::OK)
        .header("Content-Type", "text/plain")
        .body(Body::from("Multipart data processed successfully"))
        .unwrap()
}
