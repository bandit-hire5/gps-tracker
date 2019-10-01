FROM golang:1.12.1-stretch

WORKDIR $GOPATH/src/github.com/gps/gps-tracker/

COPY go.mod go.sum ./

ENV GO111MODULE=on
RUN go mod vendor

COPY . ./

RUN go build -o ./bin/main -v ./main.go

RUN cd ./
RUN ls -la