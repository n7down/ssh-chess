FROM ubuntu

RUN apt update
RUN apt install git golang net-tools iputils-ping -y

ENV PATH=/usr/lib/go/bin:$PATH
ENV GOPATH=/go

COPY app/* $GOPATH/src/github.com/n7down/ssh-chess/app/
COPY common/* $GOPATH/src/github.com/n7down/ssh-chess/common/
COPY logger/* $GOPATH/src/github.com/n7down/ssh-chess/logger/

WORKDIR $GOPATH/src/github.com/n7down/ssh-chess/app

RUN go get -v ./...
RUN CGO_ENABLED=0 go build -o /usr/bin/ssh-chess *.go

WORKDIR $GOPATH/bin
RUN apt install -y openssh-client && ssh-keygen -t rsa -N "" -f id_rsa

CMD ["ssh-chess", "-d"]

#FROM golang:latest as builder 
#WORKDIR $GOPATH/src/github.com/n7down/ssh-chess

#COPY app/* $GOPATH/src/github.com/n7down/ssh-chess/app/
#COPY common/* $GOPATH/src/github.com/n7down/ssh-chess/common/

#RUN go get -v ./... && CGO_ENABLED=0 go build -o /usr/bin/ssh-chess app/*.go

#FROM alpine:latest
#COPY --from=builder /usr/bin/ssh-chess /usr/bin/
#RUN apk add --update openssh-client && ssh-keygen -t rsa -N "" -f id_rsa
#ENTRYPOINT ssh-chess 

