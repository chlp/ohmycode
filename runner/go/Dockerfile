# docker-compose up --build --remove-orphans --force-recreate

FROM golang:latest
WORKDIR /app
COPY runner.sh /app
RUN chmod +x runner.sh
ENTRYPOINT ["./runner.sh"]
