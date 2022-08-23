FROM node:16-alpine
WORKDIR /app
COPY package*.json /
RUN npm ci
COPY src src
CMD node src/auth.js