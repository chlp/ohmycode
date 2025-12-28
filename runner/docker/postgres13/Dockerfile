FROM postgres:13.15
RUN apt-get update && apt-get install -y --no-install-recommends inotify-tools \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY runner.sh /app/runner.sh
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/runner.sh /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["postgres"]
