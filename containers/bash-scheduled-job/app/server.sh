#!/bin/bash

# Signal handling to make sure the script terminates properly
function cleanup() {
  rm -f response;
  echo "Scheduled job exited at $(date).";
}
trap cleanup INT TERM EXIT

# In this example, we use named pipe as buffers to achieve dynamic response from the server
# Named pipes employ a FIFO communication channel so first we create the response FIFO
rm -f response
mkfifo response

# This function handles the scheduled job (triggered by a POST request to /)
function handle_scheduled_job() {
  RESPONSE=$(echo -e 'HTTP/1.1 200 OK\r\n');
  bash script.sh;
}

# This function handles all other requests
function handle_non_post_request() {
  RESPONSE=$(echo -e 'HTTP/1.1 200 OK\r\n');
}

# Process the request, parse it and route it to a response handler
function handle_request() {

  # Read request, parse each line and break until we find an empty line
  while read line; do
    echo "$line"
    trline=$(echo "$line" | tr -d '[\r\n]')

    [ -z "$trline" ] && break

    # Regex to detect HTTP headline
    HEADLINE_REGEX='(.*?)\s(.*?)\sHTTP.*?'
    # Check if the line matches the headline regex
    [[ "$trline" =~ $HEADLINE_REGEX ]] &&
      # If this is an headline, we save the HTTP method and path in a variable
      REQUEST=$(echo "$trline" | sed -E "s/$HEADLINE_REGEX/\1 \2/")
  done

  # Route to the response handler based on the HTTP method and path
  case "$REQUEST" in
    "POST /")       handle_scheduled_job ;;
    *)              handle_non_post_request ;;
  esac

  # Send the response to the named pipe
  echo -e "$RESPONSE" > response

}

while true; do
  cat response | nc -q 0 -l 8080 | handle_request
done
