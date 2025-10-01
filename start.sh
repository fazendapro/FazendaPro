#!/bin/sh

export PORT=${PORT:-8080}

sed -i "s/8080/$PORT/g" /etc/nginx/nginx.conf

nginx -g "daemon off;"
