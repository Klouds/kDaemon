ADD kdaemon:/kdaemon
ADD env2conf.sh:/env2conf.sh
RUN env2conf.sh
EXPOSE 1337 13337
ENTRYPOINT ["/kdaemon/kdaemon"]
