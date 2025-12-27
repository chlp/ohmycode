FROM alpine:latest
RUN apk add --no-cache bash pandoc inotify-tools
WORKDIR /app
COPY runner.sh /app
RUN chmod +x /app/runner.sh
ENTRYPOINT ["./runner.sh"]
