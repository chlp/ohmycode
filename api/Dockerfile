# docker-compose up --build --remove-orphans

FROM php:8.2
RUN docker-php-ext-install mysqli
WORKDIR /app
EXPOSE 80
CMD ["php", "-S", "0.0.0.0:80", "-t", "/app/public"]
