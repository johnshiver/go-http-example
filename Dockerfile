FROM golang:1.13-stretch

WORKDIR /go/src/asapp_challenge
COPY . .

RUN go get -d -v ./...
RUN go build -v -o asapp cmd/httpd/*

ENTRYPOINT /go/src/asapp_challenge/asapp

EXPOSE 8080