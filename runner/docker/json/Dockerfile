FROM alpine:latest
RUN apk add --no-cache bash jq inotify-tools
WORKDIR /app
COPY runner.sh /app
RUN chmod +x /app/runner.sh
ENTRYPOINT ["./runner.sh"]
