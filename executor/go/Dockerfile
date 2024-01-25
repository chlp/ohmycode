# docker-compose up --build --remove-orphans --force-recreate

FROM golang:latest
WORKDIR /app
COPY executor.sh /app
RUN chmod +x executor.sh
ENTRYPOINT ["./executor.sh"]
