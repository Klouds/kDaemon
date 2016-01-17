FROM golang:latest
WORKDIR /kdaemon
RUN go get https://github.com/superordinate/kdaemon/
RUN go build https://github.com/superordinate/kdaemon/
RUN /kdaemon/env2conf.sh
EXPOSE 1337 13337
ENTRYPOINT ["/kdaemon/kdaemon"]
