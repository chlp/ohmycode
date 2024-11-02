# OhMyCode

* Keep notes
* Write and execute code
* Share with ease
* Collaborate in real-time
* Lightning-fast
* Use public server or deploy in your private network

Try it -> https://ohmycode.work/

![OhMyCode preview](OhMyCode-preview.png)

# Build and run

Run all services locally together (from the root of the repository):
```bash
cd api/docker && docker compose up -d --build --remove-orphans --force-recreate &&  cd../../  && \
cd ../../runner-go/docker && docker compose up -d --build --remove-orphans --force-recreate && cd../../
```
open http://localhost:52674/

---

Or configure and run separately:
1. api _(could serve client files)_:
    * `cd api`
    * `cp api-conf-example.json api-conf.json` and fill
    * `cd docker && docker compose up --build --remove-orphans --force-recreate`\
      or `cd cmd` & `GOOS=linux GOARCH=amd64 go build -o ohmycode_api` and run binary
2. runner:
    * `cd runner`
    * `cp conf-example.json conf.json` and fill
    * `cd docker && docker compose up --build --remove-orphans --force-recreate`\
      or `cd cmd` & `GOOS=linux GOARCH=amd64 go build -o ohmycode_runner` and run binary
3. client _(if you want to serve it separately)_:
    * `cd client`
    * `cp public/js/conf-example.json public/js/conf.json` and fill it with full URL to the api
    * `cd docker && docker compose up --build --remove-orphans --force-recreate`
