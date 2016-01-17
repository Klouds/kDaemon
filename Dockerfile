FROM golang:latest
WORKDIR /kdaemon
RUN go get github.com/superordinate/kdaemon/
RUN go build github.com/superordinate/kdaemon/
RUN $GOPATH/src/kdaemon/env2conf.sh
EXPOSE 1337 13337
ENTRYPOINT ["/kdaemon/kdaemon"]
