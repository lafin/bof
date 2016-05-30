## Best of <something>
___

```
$ docker run -d --restart=always -p 27017:27017 --name mongo -d mongo
$ git clone https://github.com/lafin/bof.git
$ cd bof
$ docker build -t bof .
$ docker run -e "CLIENT_ID=<...>" -e "CLIENT_EMAIL=<...>" -e "CLIENT_PASSWORD=<...>" -e "DB_SERVER=mongo" --link mongo:mongo --rm bof
```
