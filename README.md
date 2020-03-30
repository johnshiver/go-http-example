## Installation

Docker + docker-compose must be installed to run the chat server + database.
https://docs.docker.com/compose/install/

Gomock must be installed to run tests.

To install gomock (assuming go is installed)
```shell script
GO111MODULE=on go get github.com/golang/mock/mockgen@latest
```

## Build and Run chat server

Using Makefile:
```shell script
make run
```
Usually the server restarts a few times before the db comes up.
Once the logs stop writing the server should be ready to accept requests
on localhost port 8080.

## Run tests

Once gomock is installed, using Makefile:
```shell script
make test
```

With race check:
```shell script
make test-race
```

## Chat cli / integration testing
Once server is up and running you can run some quick tests to check
that everything is working using the chat cli.

To create cli binary (assuming go is installed)
```shell script
go build -o chatCli cmd/cli/*
```

Create test users:
```shell script
./chatCli -cmd=createUsers
```

Send messages:
```shell script
./chatCli -user=yoda -rec=luke -cmd=sendMsg
```

Get messages:
```shell script
./chatCli -user=yoda -rec=luke -cmd=getMsg -start=0 -limit=5 | jq '.'
```
