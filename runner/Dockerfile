FROM php:8.2-cli-alpine
RUN apk add --no-cache bash
WORKDIR /app
COPY conf.json /app
COPY conf-example.json /app
COPY tools.php /app
COPY requestsReceiver.php /app
COPY resultsSender.php /app
ENTRYPOINT ["sh", "-c", "php requestsReceiver.php & php resultsSender.php"]
