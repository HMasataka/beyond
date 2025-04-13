FROM node:20-alpine

WORKDIR /opt/firebase

RUN apk update \
    && apk --no-cache add openjdk17-jre-headless \
    && rm -rf /var/cache/apk/*

RUN npm install -g firebase-tools
