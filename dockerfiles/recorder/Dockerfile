FROM ubuntu

RUN apt update
RUN apt install git golang net-tools iputils-ping -y

ENV PATH=/usr/lib/go/bin:$PATH
ENV GOPATH=/go

COPY services/recorder/* $GOPATH/src/github.com/n7down/ssh-chess/services/recorder/
COPY services/db/* $GOPATH/src/github.com/n7down/ssh-chess/services/db/
COPY common/* $GOPATH/src/github.com/n7down/ssh-chess/common/
COPY logger/* $GOPATH/src/github.com/n7down/ssh-chess/logger/

COPY wait-for-it.sh /usr/bin/
RUN chmod +x /usr/bin/wait-for-it.sh

WORKDIR $GOPATH/src/github.com/n7down/ssh-chess/services/recorder

RUN go get -v ./...
RUN CGO_ENABLED=0 go build -o /usr/bin/recorder *.go

WORKDIR $GOPATH/bin

#CMD ["recorder"]

#FROM golang:latest as builder

#WORKDIR $GOPATH/src/github.com/n7down/ssh-chess/services/recorder
#ADD services/recorder/ . 

#RUN go get -v ./...
#RUN CGO_ENABLED=0 go build -o /usr/bin/recorder *.go

#FROM alpine:latest
#COPY --from=builder /usr/bin/recorder /usr/bin/
#ENTRYPOINT recorder 

