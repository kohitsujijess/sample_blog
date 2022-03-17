FROM golang:1.17-alpine
WORKDIR /app
RUN apk add gcc musl-dev
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
CMD ["go", "run", "main.go"]
