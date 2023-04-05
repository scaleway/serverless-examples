<?php

require __DIR__ . "/vendor/autoload.php";

function run($event, $context)
{
    // Get values from function environment variables
    $s3_access_key = getenv('S3_ACCESS_KEY');
    $s3_secret = getenv('S3_SECRET_KEY');
    $s3_endpoint = getenv('S3_ENDPOINT');
    $s3_region = getenv('S3_REGION');
    $s3_bucket = getenv('S3_BUCKET');

    $s3_key = 'example-key';

    // Connect to S3
    $s3 = new Aws\S3\S3Client([
        'region'      => $s3_region,
        'version'     => 'latest',
        'endpoint'    => $s3_endpoint,
        'credentials' => [
            'key'     => $s3_access_key,
            'secret'  => $s3_secret,
        ]
    ]);

    // Write to the key
    $result = $s3->putObject([
        'Bucket' => $s3_bucket,
        'Key'    => $s3_key,
        'Body'   => 'This is from the PHP example',
    ]);

    // Output the result of the S3 API operation
    print_r($result);

    return [
        'body' => 'S3 function succeeded',
    ];
}

?>
