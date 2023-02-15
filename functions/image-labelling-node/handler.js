import * as tf from '@tensorflow/tfjs';
import * as mobilenet from '@tensorflow-models/mobilenet';
import * as imagedata from '@andreekeberg/imagedata';
import * as path from 'path';
import { fileURLToPath } from 'url';
import S3 from "aws-sdk/clients/s3.js";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const s3EndpointUrl = process.env.S3_ENDPOINT_URL;
const userAccessKey = process.env.USER_ACCESS_KEY;
const userSecretKey = process.env.USER_SECRET_KEY;

async function handle(event, context, callback){

  // useful for function logging
  console.log(event)
  console.log(context)

  const sourceBucket = event.queryStringParameters.sourceBucket;
  const sourceKey = event.queryStringParameters.sourceKey;

  const response = {
    statusCode: 200,
    headers: {"Content-Type": ["application/json"]},
    body: {
      labels: await classifyImage(imageToTensor(await getImageFromS3(sourceBucket, sourceKey, s3EndpointUrl, userAccessKey, userSecretKey)))
    }
  };

  return response;

};


async function getImageFromS3(sourceBucket, sourceKey, s3EndpointUrl, userAccessKey, userSecretKey){

  const s3 = new S3({
    endpoint: s3EndpointUrl,
    credentials: {
      accessKeyId: userAccessKey,
      secretAccessKey: userSecretKey,
    },
  });

  const typeMatch = sourceKey.match(/\.([^.]*)$/);
  if (!typeMatch) {
    console.log("Could not determine the image type.");
  }

  const imageType = typeMatch[1].toLowerCase();
  if (["jpeg", "jpg", "png"].includes(imageType) !== true) {
    console.log(`Unsupported image type: ${imageType}`);
  }

  try {
    const params = {
        Bucket: sourceBucket,
        Key: sourceKey,
    };
    var imageObject = await s3.getObject(params).promise();
  } catch (error) {
    console.log(error);
  }

  const image = imagedata.getSync(imageObject["Body"])

  return image
}

function imageToTensor(image){

  const numChannels = 3;
  const numPixels = image.width * image.height;
  const values = new Int32Array(numPixels * numChannels);
  let pixels = image.data

  for (let i = 0; i < numPixels; i++) {
      for (let channel = 0; channel < numChannels; ++channel) {
          values[i * numChannels + channel] = pixels[i * 4 + channel];
      }
  }

  const outShape = [image.height, image.width, numChannels];
  const imageTensor = tf.tensor3d(values, outShape, 'int32');

  return imageTensor
}

async function classifyImage(imageTensor){

  const model = await mobilenet.load();
  const predictions = await model.classify(imageTensor);

  return predictions
}

export {handle}