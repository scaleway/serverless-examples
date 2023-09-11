<?php

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
        // If the status code is not in the 2XX range, the message is considered
        // failed and is retried. In total, there are 3 retries.
        "statusCode" => 200,
        "headers" => ["Content-Type" => "text/plain"],
        // Because triggers are asynchronous, the response body is ignored.
	    // It's kept here when testing locally.
        "body" => $result,
    ];
}
