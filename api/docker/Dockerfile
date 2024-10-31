FROM golang:1.21-alpine as builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o api_app cmd/main.go

FROM alpine:latest
COPY --from=builder /app/api_app /app/api-conf*.json /app/
WORKDIR /app
ENTRYPOINT ["/app/api_app"]
