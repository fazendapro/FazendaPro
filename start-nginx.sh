# Arquivo: start-nginx.sh

#!/bin/sh
# Substitui a variável no template e salva o resultado no arquivo de configuração principal do Nginx.
envsubst '${PORT}' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf
# Inicia o Nginx em primeiro plano
nginx -g 'daemon off;'