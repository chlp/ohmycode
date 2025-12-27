FROM php:8.2-cli-alpine
RUN apk add --no-cache bash inotify-tools && docker-php-ext-install pcntl
WORKDIR /app
COPY runner.sh /app
RUN chmod +x /app/runner.sh
ENTRYPOINT ["./runner.sh"]
