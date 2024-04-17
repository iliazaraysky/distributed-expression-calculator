FROM node:14-alpine

WORKDIR /app

COPY package*.json ./

RUN npm install
RUN npm install jwt-decode

COPY . .

EXPOSE 3000

CMD ["npm", "start"]