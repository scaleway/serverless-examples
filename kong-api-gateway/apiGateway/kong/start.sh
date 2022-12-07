#!/sbin/tini /bin/sh
echo "Parsing kong.yml template"
envsubst < /kong.yml.template > /kong.yml
echo "Starting Kong API Gateway DB-less"
/usr/local/bin/kong start -v -c /kong.conf
