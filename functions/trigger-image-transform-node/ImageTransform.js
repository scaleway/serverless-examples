// Add dependencies
console.log("init");
const { S3Client, PutObjectCommand, GetObjectCommand } = require("@aws-sdk/client-s3");
const sharp = require("sharp");

// Get connexion information from secret environment variables
const S3_ACCESS_KEY=process.env.S3_ACCESS_KEY;
const S3_ACCESS_KEY_ID=process.env.S3_ACCESS_KEY_ID;
const SOURCE_BUCKET=process.env.SOURCE_BUCKET;
const DESTINATION_BUCKET=process.env.DESTINATION_BUCKET;
const S3_REGION=process.env.S3_REGION;
const RESIZED_WIDTH=process.env.RESIZED_WIDTH;
const S3_ENDPOINT = `https://s3.${S3_REGION}.scw.cloud`;

let width = parseInt(RESIZED_WIDTH, 10);
if (width < 1 || width > 1000) {
  width = 200;
}

// Create S3 service object
const s3Client = new S3Client({
  credentials: {
    accessKeyId: S3_ACCESS_KEY_ID,
    secretAccessKey: S3_ACCESS_KEY,
  },
  endpoint: S3_ENDPOINT,
  region: S3_REGION
});

// Handler
exports.handle = async (event, context, callback) => {
  // Read options from the event parameter.
  console.log(event);
  // Object key may have spaces or unicode non-ASCII characters.
  const srcKey = event.body;
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
    const input = {
      Bucket: SOURCE_BUCKET,
      Key: srcKey,
    };
    //var origimage = await s3.getObject(params).promise();
    const getObjectCommand = new GetObjectCommand(input);
    var getObjectResult = await s3Client.send(getObjectCommand);
  } catch (error) {
    console.error(error);
    return error;
  }

  // Use the sharp module to resize the image.
  try {
    var sharpImg = sharp().resize({ width, withoutEnlargement: true })
    getObjectResult.Body.pipe(sharpImg);
  } catch (error) {
    console.error(error);
    return error;
  }

  // Upload the image as a Buffer to the destination bucket
  try {
    const destinput = {
      Bucket: DESTINATION_BUCKET,
      Key: dstKey,
      Body: await sharpImg.toBuffer(),
      ContentType: "image",
    };
    //await s3.putObject(destparams).promise();
    const putObjectCommand = new PutObjectCommand(destinput);
    const putimage = await s3Client.send(putObjectCommand);
    console.log(putimage.VersionId)
  } catch (error) {
    console.log(error);
    return error;
  }
  console.log(
    "Successfully resized " +
    SOURCE_BUCKET +
    "/" +
    srcKey +
    " and uploaded to " +
    DESTINATION_BUCKET +
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
        DESTINATION_BUCKET,
    }),
    headers: {
      "Content-Type": "application/json",
    },
  };
};


/* This is used to test locally and will not be executed on Scaleway Functions */
if (process.env.NODE_ENV === 'test') {
  import("@scaleway/serverless-functions").then(scw_fnc_node => {
    scw_fnc_node.serveHandler(exports.handle, 8081);
  });
}
