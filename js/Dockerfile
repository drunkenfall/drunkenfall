FROM node:11-alpine

COPY package*.json ./
RUN npm install --only-production

# Copy mostly static js stuff
COPY build/ ./build/
COPY config/ ./config/
COPY .eslintrc.js ./.eslintrc.js
COPY .babelrc ./.babelrc
COPY static/ ./static/
COPY index.html ./

COPY src ./src/

RUN npm run build
