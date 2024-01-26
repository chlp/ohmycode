FROM mysql:8
WORKDIR /app
COPY runner.sh /app
RUN chmod +x runner.sh
ENTRYPOINT ["./runner.sh"]
