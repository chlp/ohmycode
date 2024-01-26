# docker-compose up --build --remove-orphans --force-recreate

FROM mysql:8
WORKDIR /app
COPY runner.sh /app
RUN chmod +x runner.sh
ENTRYPOINT ["./runner.sh"]
