FROM golang:alpine

RUN ls bin
RUN ls src

ADD config config

ADD dist/qf_linux_amd64/qf /qf

ENTRYPOINT ["/qf"]