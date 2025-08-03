FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine

RUN apk add --no-cache gettext

COPY --from=builder /app/dist /usr/share/nginx/html

COPY nginx.conf.template /etc/nginx/nginx.conf.template
COPY start-nginx.sh /start-nginx.sh

RUN chmod +x /start-nginx.sh

ENV PORT 8080
EXPOSE 8080

CMD ["/start-nginx.sh"]