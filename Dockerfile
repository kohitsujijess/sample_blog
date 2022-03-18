# syntax=docker/dockerfile:1
FROM golang:1.18-alpine
WORKDIR /app
RUN apk add gcc musl-dev git
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
