FROM golang:1.21-alpine as builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o app cmd/main.go

FROM alpine:latest
COPY --from=builder /app/app /app/runner-conf*.json /app/
WORKDIR /app
ENTRYPOINT ["/app/app"]
