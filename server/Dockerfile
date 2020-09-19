FROM node:current-alpine

COPY package*.json ./

RUN npm ci --only=production

COPY *.js ./

CMD ["node", "index.js"]