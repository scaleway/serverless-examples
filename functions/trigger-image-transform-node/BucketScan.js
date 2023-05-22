
console.log("init");
const { S3Client, ListObjectsV2Command } = require("@aws-sdk/client-s3");
const { SQSClient, SendMessageCommand } = require("@aws-sdk/client-sqs");

// Get info from env variables
const S3_ACCESS_KEY = process.env.S3_ACCESS_KEY;
const S3_ACCESS_KEY_ID = process.env.S3_ACCESS_KEY_ID;
const SOURCE_BUCKET = process.env.SOURCE_BUCKET;
const S3_REGION = process.env.S3_REGION;
const SQS_ACCESS_KEY = process.env.SQS_ACCESS_KEY;
const SQS_ACCESS_KEY_ID = process.env.SQS_ACCESS_KEY_ID;
const QUEUE_URL = process.env.QUEUE_URL;
const SQS_ENDPOINT = process.env.SQS_ENDPOINT;
const S3_ENDPOINT = `https://s3.${S3_REGION}.scw.cloud`;

// Create S3 service object
const s3Client = new S3Client({
  credentials: {
    accessKeyId: S3_ACCESS_KEY_ID,
    secretAccessKey: S3_ACCESS_KEY,
  },
  endpoint: S3_ENDPOINT,
  region: S3_REGION
});

// Configure parameters for the listObjectsV2 Command
const input = {
  "Bucket": SOURCE_BUCKET
};

// Create SQS service
var sqsClient = new SQSClient({
  credentials: {
    accessKeyId: SQS_ACCESS_KEY_ID,
    secretAccessKey: SQS_ACCESS_KEY
  },
  region: "par",
  endpoint: SQS_ENDPOINT,
})

console.log("init Ok")

exports.handle = async (event, context, callback) => {
  const s3ListCommand = new ListObjectsV2Command(input);
  const s3List = await s3Client.send(s3ListCommand);
  let counter = 0;
  const contents = s3List.Contents;
  const total = contents.length;
  contents.forEach(async function (content) {
    // Infer the image type from the file suffix.
    const srcKey = content.Key;
    const typeMatch = srcKey.match(/\.([^.]*)$/);
    if (!typeMatch) {
      console.error("Could not determine the image type.")
      return {
        statusCode: 500,
        body:"Could not determine the image type."
      };
    }
    const imageType = typeMatch[1].toLowerCase();
    // Check that the image type is supported
    if (!["jpeg", "jpg", "png"].includes(imageType)) {
      console.error(`Unsupported image type: ${imageType}`)
      return {
        statusCode: 500,
        body:`Unsupported image type: ${imageType}`
      };
    }
    else {
      try {
        var sendMessageCommand = new SendMessageCommand({
          QueueUrl: QUEUE_URL,
          MessageBody: srcKey,
        });
        var sendMessage = await sqsClient.send(sendMessageCommand);
        console.log(sendMessage.MessageId);
        counter += 1;
      } catch (error) {
        console.error(error);
      }
    }
  });
  return {
    statusCode: 200,
    body: JSON.stringify({
      status: "Done",
      message: `All images from the bucket have been processed over ${total} files in the bucket`,
    }),
    headers: {
      "Content-Type": "application/json",
    },
  };
};



/* This is used to test locally and will not be executed on Scaleway Functions */
if (process.env.NODE_ENV === 'test') {
  import("@scaleway/serverless-functions").then(scw_fnc_node => {
    scw_fnc_node.serveHandler(exports.handle, 8080);
  });
}
