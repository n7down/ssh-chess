FROM golang:latest as builder 
WORKDIR $GOPATH/src/github.com/n7down/ssh-chess

COPY cmd/ $GOPATH/src/github.com/n7down/ssh-chess/cmd
COPY internal/ $GOPATH/src/github.com/n7down/ssh-chess/internal

RUN go get -v -d ./... 
RUN CGO_ENABLED=0 go build -o /usr/bin/ssh-chess cmd/ssh-chess/main.go

FROM alpine:latest
COPY --from=builder /usr/bin/ssh-chess /usr/bin/
RUN apk add --update openssh-client && ssh-keygen -t rsa -N "" -f id_rsa
CMD ["ssh-chess"]
