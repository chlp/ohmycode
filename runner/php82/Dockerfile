# docker-compose up --build --remove-orphans --force-recreate

FROM php:8.2-cli
WORKDIR /app
COPY runner.sh /app
RUN chmod +x runner.sh
ENTRYPOINT ["./runner.sh"]
