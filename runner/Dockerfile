FROM php:8.2-cli
WORKDIR /app
COPY conf.json /app
COPY manager.php /app
ENTRYPOINT ["php", "manager.php"]
