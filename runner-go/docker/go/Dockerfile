FROM golang:1.21-alpine
RUN apk add --no-cache bash
WORKDIR /app
COPY runner.sh /app
RUN chmod +x /app/runner.sh
ENTRYPOINT ["./runner.sh"]
