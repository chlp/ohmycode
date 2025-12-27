FROM postgres:13.15
RUN apt-get update && apt-get install -y --no-install-recommends inotify-tools \
  && rm -rf /var/lib/apt/lists/*
