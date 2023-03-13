FROM alpine
WORKDIR /usr/src/app

ENV PORT 8080
EXPOSE 8080

RUN apk update && apk add --no-cache netcat-openbsd bash
COPY server.sh .
COPY script.sh .

CMD ["/bin/bash", "./server.sh"]
