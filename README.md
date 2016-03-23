# Quality Farm

(PoC created during 24h hackaton)

Static code analysis tool for Go. It runs tests coverage and dozen of linters, calculates total quality score and technical debt.

## Installation

### Quick Run

All services can be started at once with two commands:

```bash
docker-compose build
docker-compose up
```

Quality Farm front-end should be now available at http://docker:9000/. It is recommended to add /etc/hosts alias for docker, eg.:

```bash
192.168.99.100 docker
# or
127.0.0.1 docker localhost
```

## Development

### Front-end

While working on front-end part, you can start back-end services with `docker-compose up redis websocket server`. Now you can start front-end separately:

```bash
cd webapp/
npm install typings webpack-dev-server rimraf webpack -g
npm install
npm start
```

App should be available at http://localhost:3000/

### Back-end

While working on API server, you can start other services with `docker-compose up redis websocket`, then start front-end as described above and build & run server part:

```bash
go install ./cmd/server/ && server
```
