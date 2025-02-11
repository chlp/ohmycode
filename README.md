# OhMyCode

* Keep notes
* Write and execute code
* Share with ease
* Collaborate in real-time
* Lightning-fast
* Use public server or deploy in your private network

Try it -> https://ohmycode.work/

![OhMyCode preview](OhMyCode-preview.png)

Technical Article: ["OhMyCode: System Design Reflections"](https://chlp8.medium.com/ohmycode-system-design-reflections-07c26f91b861)

# Build and run

Run all services locally together (from the root of the repository):
```bash
docker compose -f api/docker/docker-compose.yml up -d --build --remove-orphans --force-recreate && \
docker compose -f runner/docker/docker-compose.yml up -d --build --remove-orphans --force-recreate
```
open http://localhost:52674/

---

Or configure and run separately:
1. api & client:
    * `cd api`
    * `cp api-conf-example.json api-conf.json` and fill
    * `cd docker && docker compose up --build --remove-orphans --force-recreate`\
      or `GOOS=linux GOARCH=amd64 go build -o ohmycode_api cmd/main.go` and run binary
2. runner:
    * `cd runner`
    * `cp conf-example.json conf.json` and fill
    * `cd docker && docker compose up --build --remove-orphans --force-recreate`\
      or `GOOS=linux GOARCH=amd64 go build -o ohmycode_runner cmd/main.go` and run binary
