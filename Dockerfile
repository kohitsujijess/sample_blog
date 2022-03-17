FROM golang:1.17-alpine
WORKDIR /app
RUN apk add gcc musl-dev
COPY . .
# RUN go mod download
