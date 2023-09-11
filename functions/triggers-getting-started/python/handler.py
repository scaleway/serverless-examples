from typing import TYPE_CHECKING
import functools
import operator
import http

if TYPE_CHECKING:
    from scaleway_functions_python.framework.v1.hints import Context, Event, Response


def factorial(n: int) -> int:  # pylint: disable=invalid-name
    """Return the factorial of n >= 0."""
    return functools.reduce(operator.mul, range(1, n + 1), 1)


def handler(event: "Event", _context: "Context") -> "Response":
    """Compute factorial of the number passed in the trigger message."""

    if event["httpMethod"] != "POST":
        # SQS triggers are sent as POST requests.
        status = http.HTTPStatus.METHOD_NOT_ALLOWED
        return {
            "headers": {"Content-Type": "text/plain"},
            "statusCode": status.value,
            "body": status.description,
        }

    # The content of the SQS message is passed in the body.
    n = int(event["body"])  # pylint: disable=invalid-name
    result = factorial(n)

    print(f"python: factorial of {n} is {result}")

    return {
        "headers": {"Content-Type": "text/plain"},
        # If the status code is not in the 2XX range, the message is considered
        # failed and is retried. In total, there are 3 retries.
        "statusCode": http.HTTPStatus.OK.value,
        # Because triggers are asynchronous, the response body is ignored.
        # It's kept here when testing locally.
        "body": str(result),
    }


if __name__ == "__main__":
    from scaleway_functions_python import local

    # Example usage:
    # curl -X POST -d 5 http://localhost:8080
    local.serve_handler(handler)
