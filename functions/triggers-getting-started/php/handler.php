<?php

require __DIR__ . '/vendor/autoload.php';

function factorial($n)
{
    $result = 1;
    for ($i = 1; $i <= $n; $i++) {
        $result *= $i;
    }
    return $result;
}

function handle($event, $context)
{
    if ($event["httpMethod"] !== "POST") {
        # SQS triggers are sent as POST requests.
        return [
            "statusCode" => 405,
            "headers" => ["Content-Type" => "text/plain"],
            "body" => "Method Not Allowed",
        ];
    }

    # The content of the SQS message is passed in the body.
    $n = intval($event["body"]);
    $result = factorial($n);

    echo "php: factorial of $n is $result\n";

    return [
        "statusCode" => 200,
        "headers" => ["Content-Type" => "text/plain"],
        "body" => $result,
    ];
}
