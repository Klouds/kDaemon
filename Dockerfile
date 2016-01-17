FROM alpine:3.1
EXPOSE 1337 13337
ENTRYPOINT ["/bin/sh"]
CMD ["env2conf.sh"]
