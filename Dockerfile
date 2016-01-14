FROM golang:1.5.3-onbuild
EXPOSE 1337 13337
ENTRYPOINT [/go/bin/kdaemon]
