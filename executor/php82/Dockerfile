# docker-compose up --build --remove-orphans --force-recreate

FROM php:8.2-cli
WORKDIR /app
COPY executor.sh /app
RUN chmod +x executor.sh
ENTRYPOINT ["./executor.sh"]
