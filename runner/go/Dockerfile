FROM golang:latest
WORKDIR /app
COPY runner.sh /app
RUN chmod +x runner.sh
ENTRYPOINT ["./runner.sh"]
