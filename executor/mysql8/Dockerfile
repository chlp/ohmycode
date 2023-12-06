# docker-compose up --build --remove-orphans --force-recreate

FROM mysql:8
WORKDIR /app
COPY executor.sh /app
RUN chmod +x executor.sh
ENTRYPOINT ["./executor.sh"]
