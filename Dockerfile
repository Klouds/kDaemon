FROM golang:latest
WORKDIR /kdaemon
RUN go get github.com/superordinate/kdaemon/
RUN go build github.com/superordinate/kdaemon/
EXPOSE 1337 13337
ENTRYPOINT ["$GOPATH/src/github.com/superordinate/kdaemon/env2conf.sh"]
