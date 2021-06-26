FROM golang:alpine

RUN ls

ADD config config

ADD dist/qf_linux_amd64/qf /qf

ENTRYPOINT ["/qf"]