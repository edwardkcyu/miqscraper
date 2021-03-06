FROM golang:1.16 as BUILDER
WORKDIR /app
COPY go.mod go.sum
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main

FROM alpine:3.12
WORKDIR /app
COPY --from=BUILDER /app/main main
ENTRYPOINT ["./main"]