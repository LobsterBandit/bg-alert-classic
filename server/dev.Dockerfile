FROM node:current-buster

WORKDIR /usr/src

COPY package*.json ./

RUN npm install

CMD ["node_modules/.bin/nodemon", "index.js"]