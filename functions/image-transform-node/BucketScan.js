/*
Work in progress : il s'agit de la fonction permettant de collecter l'ensemble des fichiers images présents 
sur un bucket (voir si on doit mettre des limites)
et d'envoyer leurs identifiants à la fonction de transformation d'image
Elle doit être déployer de préférence avec un CRON
*/

const AWS = require("aws-sdk");
const https = require("https");
const { isNull } = require("util");

// get info from env variable

const SOURCE_BUCKET = process.env.SOURCE_BUCKET;
const S3_ENDPOINT_URL = process.env.S3_ENDPOINT_URL;
const ID = process.env.ACCESS_KEY_ID;
const SECRET = process.env.ACCESS_KEY;
const TRANSFORM_URL = process.env.TRANSFORM_URL;

// Create S3 service object
const s3 = new AWS.S3({
  endpoint: S3_ENDPOINT_URL,
  credentials: {
    accessKeyId: ID,
    secretAccessKey: SECRET,
  },
});
//configure parameters for the listObjectsV2 method
const params = {
  Bucket: SOURCE_BUCKET,
};

exports.handle = async (event, context, callback) => {
  s3.listObjectsV2(params, function (err, data) {
    if (err) {
      console.log(err, err.stack); // an error occurred
    } else {
      let counter = 0;
      const contents = data.Contents;
      const total = contents.length;
      contents.forEach(function (content) {
        // Infer the image type from the file suffix.
        const srcKey = content.Key;
        const typeMatch = srcKey.match(/\.([^.]*)$/);
        if (!typeMatch) {
          return console.log("Could not determine the image type.");
        }
        const imageType = typeMatch[1].toLowerCase();
        // Check that the image type is supported
        if (["jpeg", "jpg", "png"].includes(imageType) !== true) {
          return console.log(`Unsupported image type: ${imageType}`);
        }
        try {
          console.log(TRANSFORM_URL + "?key=" + srcKey);
          https.get(TRANSFORM_URL + "?key=" + srcKey);
          counter += 1;
        } catch (error) {
          console.log(error);
        }
      });
      if (data.IsTruncated) {
        params.ContinuationToken = data.NextContinuationToken;
        return {
          statusCode: 200,
          body: JSON.stringify({
            status: "Too many items",
            message: "There are more files in your bucket that we can handle",
          }),
          headers: {
            "Content-Type": "application/json",
          },
        };
      }
      return {
        statusCode: 200,
        body: JSON.stringify({
          status: "Done",
          message: `${counter} images from the bucket have been processed over ${total} files in the bucket`,
        }),
        headers: {
          "Content-Type": "application/json",
        },
      };
    }
  });
};
