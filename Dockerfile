FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine

RUN apk add --no-cache gettext

COPY --from=builder /app/dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/nginx.conf.template
COPY start-nginx.sh /start-nginx.sh

RUN chmod +x /start-nginx.sh

# Garante que o diretório de logs existe
RUN mkdir -p /var/log/nginx

# Define a porta padrão
ENV PORT 8080
EXPOSE 8080

# Usa o script de inicialização
CMD ["/start-nginx.sh"]