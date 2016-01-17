FROM alpine:3.1
EXPOSE 1337 13337
ADD kdaemon /kdaemon/kdaemon
ADD views/ /kdaemon/views
ADD public /kdaemon/public
ADD env2conf.sh /kdaemon/env2conf.sh
ENTRYPOINT ["/bin/sh"]
CMD ["env2conf.sh"]
