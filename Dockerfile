FROM l3iggs/archlinux:latest
RUN pacman --noconfirm -Sy archlinux-keyring
RUN pacman-db-upgrade
RUN pacman --noconfirm -Syy
RUN pacman --noconfirm -Syu
RUN pacman --noconfirm -S go base-devel
RUN mkdir /go
ENV GOPATH /go
ADD . /go/src/kdaemon
WORKDIR /go/src/kdaemon
RUN go get .
RUN go build .
RUN mkdir /kdaemon
RUN cp ./config/app.conf /kdaemon/app.conf
RUN cp ./kdaemon /kdaemon/kdaemon
EXPOSE 1337 13337
ENTRYPOINT ["/kdaemon/kdaemon"]
