# REST API golang [![Build Status](http://34.123.0.188:8080/job/Go%20REST%20API/badge/icon)](http://34.123.0.188:8080/job/Go%20REST%20API/)

This is my solution to designing rest API in golang for golang school [Yalantis](https://yalantis.com/). The task description is [here](https://docs.google.com/document/d/1PPAbDVllQYpw7bFRStGB_Gbcoj2TRs9NvyPfepuAf8w/edit#heading=h.l29vobcyrk4t) 

* [api server](https://faketrello.ml/)
* [api docs(swagger)](https://faketrello.ml/docs#/)


# HOWTO build and run server

First of all, you need to set environment variables, no matter what method of starting the server you will use (see [.env.example](https://github.com/osetr/rest-api-golang/blob/master/.env.example)). Then pick preferable way.
```sh
# running with docker
docker-compose up -d --build

# running with make
migrate -path migrations -database "postgres://localhost/{db_name}?sslmode={sslmode}&user={db_user}&password={db_password}" up # to run migrations with github.com/golang-migrate/migrate
make
./v1
```

# HOWTO run tests

```sh
# runs tests
make test
```
