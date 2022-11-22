use anyhow::Result;
use dfdx::prelude::*;
use s3::{creds::Credentials, Bucket, Region};

pub const MODEL_PATH: &str = "mnist-classifier.npz";

// our network structure
pub type Mlp = (
    (Linear<784, 512>, ReLU),
    (Linear<512, 128>, ReLU),
    (Linear<128, 32>, ReLU),
    Linear<32, 10>,
);

pub fn scaleway_bucket_from_env(region: &str, bucket_name: &str) -> Result<Bucket> {
    let region = Region::Custom {
        region: region.to_string(),
        endpoint: format!("s3.{}.scw.cloud", region),
    };
    let credentials =
        Credentials::from_env_specific(Some("SCW_ACCESS_KEY"), Some("SCW_SECRET_KEY"), None, None)?;

    Ok(Bucket::new(bucket_name, region, credentials)?)
}
