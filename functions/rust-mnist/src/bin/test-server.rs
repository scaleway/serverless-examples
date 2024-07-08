use axum::{routing::post, Router};

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/", post(function::handler));
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app.into_make_service())
        .await
        .unwrap();
}
