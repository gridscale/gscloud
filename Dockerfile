FROM golang:alpine

COPY . /src

RUN cd /src && \
    apk add --no-cache make git && \
    make

FROM alpine:latest

COPY --from=0 /src/gscloud /usr/bin/gscloud
ENTRYPOINT ["/usr/bin/gscloud"]
