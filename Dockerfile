FROM debian:sid
ENV GOPATH /go
ENV PATH $PATH:/go/bin
RUN apt-get update
RUN apt-get install -y golang build-essential git
RUN mkdir /go
RUN git clone https://github.com/superordinate/kdaemon
WORKDIR /kdaemon
RUN go get https://github.com/superordinate/kdaemon
EXPOSE 1337 13337
ENTRYPOINT ["/bin/bash"]
CMD ["/kdaemon/env2conf.sh"]
