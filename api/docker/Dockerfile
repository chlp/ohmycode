FROM golang:1.21-alpine as builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o api_app cmd/main.go

FROM alpine:latest
COPY --from=builder /app/api_app /app/
COPY --from=builder /app/api-conf.prod.json /app/api-conf.json
WORKDIR /app
ENTRYPOINT ["/app/api_app"]
