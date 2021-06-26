FROM golang:alpine

COPY qf /qf

ENTRYPOINT ["/qf"]