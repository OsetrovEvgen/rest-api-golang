FROM golang:1.15.7-buster

COPY . /go/src/app

WORKDIR /go/src/app

RUN go build -v ./api/v1/

ENTRYPOINT [ "./v1" ]