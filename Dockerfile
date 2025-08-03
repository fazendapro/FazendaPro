FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine

# Instala gettext para substituição de variáveis
RUN apk add --no-cache gettext

# Copia os arquivos da aplicação
COPY --from=builder /app/dist /usr/share/nginx/html

# Copia a configuração do nginx
COPY nginx.conf /etc/nginx/nginx.conf.template

# Cria diretório de logs
RUN mkdir -p /var/log/nginx

# Define a porta padrão
ENV PORT=8080
EXPOSE 8080

# Script de inicialização inline (mais simples e confiável)
CMD ["/bin/sh", "-c", "envsubst '${PORT}' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf && nginx -g 'daemon off;'"]