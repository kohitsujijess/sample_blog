FROM golang:1.17-alpine
RUN mkdir /go/src/app
WORKDIR /go/src/app
RUN apk add --no-cache gcc musl-dev
COPY ./go.mod /go/src/app
RUN go mod tidy
CMD ["go", "run", "main.go"]