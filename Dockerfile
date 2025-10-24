FROM ubuntu:latest
LABEL authors="makskozlov"

ENTRYPOINT ["top", "-b"]