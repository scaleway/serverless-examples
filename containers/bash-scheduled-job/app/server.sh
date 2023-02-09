#!/bin/bash

### Create the response FIFO
rm -f response
mkfifo response

function handle_scheduled_job() {
  RESPONSE=$(echo -e 'HTTP/1.1 200 OK\r\n');
  bash script.sh;
}

function handle_handshake() {
  RESPONSE=$(echo -e 'HTTP/1.1 200 OK\r\n');
}

function handleRequest() {

  ## Read request, parse each line and breaks until empty line
  while read line; do
    echo "$line"
    trline=$(echo "$line" | tr -d '[\r\n]')

    [ -z "$trline" ] && break

    HEADLINE_REGEX='(.*?)\s(.*?)\sHTTP.*?'
    [[ "$trline" =~ $HEADLINE_REGEX ]] &&
      REQUEST=$(echo "$trline" | sed -E "s/$HEADLINE_REGEX/\1 \2/")
  done

  ## Route to the response handler based on the REQUEST match
  case "$REQUEST" in
    "POST /")       handle_scheduled_job ;;
    *)              handle_handshake ;;
  esac

  echo -e "$RESPONSE" > response

}

while true; do
  cat response | nc -q 0 -l 8080 | handleRequest
done