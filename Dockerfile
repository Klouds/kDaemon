FROM golang:latest
WORKDIR /kdaemon
RUN go get github.com/superordinate/kdaemon/
RUN go build github.com/superordinate/kdaemon/
RUN $GOPATH/src/github.com/superordinate/kDaemon
EXPOSE 1337 13337
ENTRYPOINT ["/kdaemon/kdaemon"]
