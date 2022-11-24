FROM kong:alpine

USER root

COPY kong.yml.template /kong.yml.template
COPY kong.conf /kong.conf
COPY start.sh /start.sh

RUN chmod a+w /kong.yml.template \
 && chmod +x /start.sh \
 && apk add --no-cache \
      tini \
      gettext \
 && rm -rf \
      /var/cache/* \
      /root/.cache/* \
      /docker-entrypoint.sh

ENV KONG_NGINX_DAEMON=off

ENTRYPOINT ["/start.sh"]
