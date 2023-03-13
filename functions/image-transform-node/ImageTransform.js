// Add dependencies
const AWS = require("aws-sdk");
const sharp = require("sharp");

// Get connexion information from secret environment variables
const srcBucket = process.env.SOURCE_BUCKET;
const S3_ENDPOINT_URL = process.env.S3_ENDPOINT_URL;
const ID = process.env.ACCESS_KEY_ID;
const SECRET = process.env.ACCESS_KEY;
const dstBucket = process.env.DESTINATION_BUCKET;

let width = parseInt(process.env.RESIZED_WIDTH, 10);
if (width < 1 || width > 1000) {
  width = 200;
}

// Create S3 service object
const s3 = new AWS.S3({
  endpoint: S3_ENDPOINT_URL,
  credentials: {
    accessKeyId: ID,
    secretAccessKey: SECRET,
  },
});

// Handler
exports.handle = async (event, context, callback) => {
  // Read options from the event parameter.
  console.log("Reading options from event");

  // Object key may have spaces or unicode non-ASCII characters.
  const srcKey = event.queryStringParameters.key;
  console.log(srcKey);
  const dstKey = "resized-" + srcKey;
  // Infer the image type from the file suffix.
  const typeMatch = srcKey.match(/\.([^.]*)$/);
  if (!typeMatch) {
    console.log("Could not determine the image type.");
    return;
  }

  // Check that the image type is supported
  const imageType = typeMatch[1].toLowerCase();
  if (["jpeg", "jpg", "png"].includes(imageType) !== true) {
    console.log(`Unsupported image type: ${imageType}`);
    return;
  }

  // Download the image from the S3 source bucket.
  try {
    const params = {
      Bucket: srcBucket,
      Key: srcKey,
    };
    var origimage = await s3.getObject(params).promise();
  } catch (error) {
    console.log(error);
    return;
  }

  // Use the sharp module to resize the image and save in a buffer.
  try {
    var buffer = await sharp(origimage.Body)
      .resize({ width, withoutEnlargement: true })
      .toBuffer();
  } catch (error) {
    console.log(error);
    return;
  }

  // Upload the image to the destination bucket
  try {
    const destparams = {
      Bucket: dstBucket,
      Key: dstKey,
      Body: buffer,
      ContentType: "image",
    };
    await s3.putObject(destparams).promise();
  } catch (error) {
    console.log(error);
    return;
  }
  console.log(
    "Successfully resized " +
    srcBucket +
    "/" +
    srcKey +
    " and uploaded to " +
    dstBucket +
    "/" +
    dstKey
  );
  return {
    statusCode: 201,
    body: JSON.stringify({
      status: "ok",
      message:
        "Image : " +
        srcKey +
        " has successfully been resized and pushed to the bucket " +
        dstBucket,
    }),
    headers: {
      "Content-Type": "application/json",
    },
  };
};
