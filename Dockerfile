ADD . /kdaemon/
RUN env2conf.sh
EXPOSE 1337 13337
ENTRYPOINT ["/kdaemon/kdaemon"]
