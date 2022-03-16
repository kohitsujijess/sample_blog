FROM golang:1.17-alpine
RUN mkdir /go/src/app
WORKDIR /go/src/app
RUN apk add --no-cache gcc musl-dev
