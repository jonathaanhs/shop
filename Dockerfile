# Builder
FROM golang:1.17.0-alpine as builder

RUN mkdir /shop
COPY  . /shop
WORKDIR /shop/app

RUN apk add --no-cache git mercurial \
    && go get -d -v\
    && apk del git mercurial 

RUN apk add build-base

EXPOSE 8080

CMD ["go", "run", "main.go"]