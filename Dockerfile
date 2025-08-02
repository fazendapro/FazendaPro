FROM node:20.19.0

WORKDIR /app

COPY package*.json ./

RUN npm install
RUN npm install -g @nestjs/cli
RUN npm install --save-dev @types/node @types/express @types/bcryptjs @types/bcrypt

RUN apt-get update && apt-get install -y netcat-openbsd

COPY . .

RUN npm run build

COPY wait-for-postgres.sh /app/wait-for-postgres.sh
RUN chmod +x /app/wait-for-postgres.sh

EXPOSE 3000 5555

CMD ["npm", "run", "start"]