### Best of ___something___ ![backend](https://github.com/lafin/bof/workflows/backend/badge.svg)
___

```bash
$ docker-compose -f infra/docker-compose.bof.yml pull
$ CLIENT_PASSWORD= DB_PASSWORD= docker-compose -f infra/docker-compose.bof.yml up
```

### go deps
```sh 
$ go mod tidy && go get -u
```