# Arquivo: start-nginx.sh

#!/bin/sh
# Substitui a variável no template e salva o resultado no arquivo de configuração padrão do Nginx.
envsubst '${PORT}' < /etc/nginx/nginx.conf.template > /etc/nginx/conf.d/default.conf
# Inicia o Nginx em primeiro plano
nginx -g 'daemon off;'