use anyhow::{Context, Result};
use serde::{Deserialize, Serialize};
use std::{env, fmt::Display, fs::File, io::Write};

use hyper::{Body, Request, Response, StatusCode};

use common::{scaleway_bucket_from_env, Mlp, MODEL_PATH};
use dfdx::prelude::*;

pub async fn load_model() -> Result<Mlp> {
    // Obtain the model from the bucket
    let bucket = scaleway_bucket_from_env(
        &env::var("SCW_DEFAULT_REGION").unwrap_or("fr-par".to_owned()),
        &env::var("S3_BUCKET")?,
    )?;
    let data = bucket.get_object(MODEL_PATH).await?;

    // Write the model to a temporary file
    let path = env::temp_dir().join(MODEL_PATH);
    let mut file = File::create(path.clone())?;
    file.write_all(data.bytes())?;

    // Load the model from the temporary file
    let mut model: Mlp = Default::default();
    model
        .load(path.clone())
        .with_context(|| format!("could not load model from {}", path.to_str().unwrap_or("")))?;
    Ok(model)
}

#[derive(Deserialize)]
struct LabelingRequest {
    data: Vec<f32>,
}

#[derive(Serialize)]
struct LabelingResponse {
    output: Vec<f32>,
}

fn handle_error<Err: Display>(err: Err, status: Option<StatusCode>) -> Response<Body> {
    let status = status.unwrap_or(StatusCode::INTERNAL_SERVER_ERROR);
    Response::builder()
        .status(status)
        .body(Body::from(format!("{}", err)))
        .unwrap()
}

pub async fn sync_handler(req: Request<Body>) -> Response<Body> {
    let (sender, receiver) = std::sync::mpsc::channel();
    std::thread::spawn(move || {
        let rt = tokio::runtime::Builder::new_current_thread()
            .enable_all()
            .build()
            .unwrap();
        rt.block_on(async {
            let resp = handler(req).await;
            sender
                .send(resp)
                .expect("Error when sending response back to rt")
        })
    });
    receiver
        .recv()
        .expect("Error when receiving response from thread")
}

pub async fn handler(req: Request<Body>) -> Response<Body> {
    let model = match load_model().await {
        Ok(m) => m,
        Err(e) => return handle_error(e, None),
    };

    let body = req.into_body();

    let res = match hyper::body::to_bytes(body).await {
        Ok(body) => serde_json::from_slice::<LabelingRequest>(&body),
        Err(e) => return handle_error(e, None),
    };

    let request = match res {
        Ok(req) if req.data.len() >= 784 => req,
        _ => return handle_error("invalid request", Some(StatusCode::BAD_REQUEST)),
    };

    let mut img = Tensor1D::<784>::zeros();
    img.mut_data().copy_from_slice(&request.data[0..784]);

    let output = model.forward(img);

    let response = LabelingResponse {
        output: output.data().to_vec(),
    };

    let body = match serde_json::to_string(&response) {
        Ok(body) => body,
        Err(e) => return handle_error(e, None),
    };

    Response::builder()
        .status(StatusCode::OK)
        .body(Body::from(body))
        .unwrap()
}

#[cfg(test)]
mod tests {
    use hyper::{Body, Request, StatusCode};
    use serde_json::json;
    use tokio_test::*;

    #[tokio::test]
    async fn test_handler() {
        let json = json!({
            "data": vec![0.0 as f32; 784],
        });

        let request = assert_ok!(Request::builder().body(Body::from(json.to_string())));

        let (parts, body) = crate::handler(request).await.into_parts();
        println!(
            "{}",
            String::from_utf8(hyper::body::to_bytes(body).await.unwrap().to_vec()).unwrap()
        );

        assert_eq!(parts.status, StatusCode::OK)
    }
}
