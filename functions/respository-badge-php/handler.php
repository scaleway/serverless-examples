<?php

require __DIR__ . "/vendor/autoload.php";

use Twig\Environment;
use Twig\Error\LoaderError;
use Twig\Error\RuntimeError;
use Twig\Error\SyntaxError;
use Twig\Loader\FilesystemLoader;

/**
 * @throws RuntimeError
 * @throws SyntaxError
 * @throws LoaderError
 */
function handler($event, $context)
{
    $text = "Scaleway";
    $status = "approved";
    $color = "green";

    if (isset($event["queryStringParameters"]["text"])) {
        $text = $event["queryStringParameters"]["text"];
    }
    if (isset($event["queryStringParameters"]["status"])) {
        $status = $event["queryStringParameters"]["status"];
    }
    if (isset($event["queryStringParameters"]["color"])) {
        $color = $event["queryStringParameters"]["color"];
    }

    $textWidth = strlen($text) * 71;
    $statusWidth = strlen($status) * 71;

    $loader = new FilesystemLoader(__DIR__ . "/views/");
    $twig = new Environment($loader, []);

    return [
        "body" => $twig->render('badge.twig', [
            "text" => [
                "text" => $text,
                "width" => $textWidth
            ],
            "status" => [
                "text" => $status,
                "width" => $statusWidth
            ],
            "color" => trim($color, "#")
        ]),
        "headers" => [
            "Content-Type" => "image/svg+xml"
        ]
    ];
}