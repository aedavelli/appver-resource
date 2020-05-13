FROM golang:alpine as builder
COPY . /go/src/github.com/aedavelli/appver-resource
ENV CGO_ENABLED 0
ENV GO111MODULE on
RUN apk add git
RUN set -e; \ 
    cd  /go/src/github.com/aedavelli/appver-resource; \ 
    for pkg in in out check; \
    do \
       go build -o /assets/$pkg github.com/aedavelli/appver-resource/$pkg; \ 
       go test -v  github.com/aedavelli/appver-resource/$pkg; \ 
    done   

FROM alpine:edge AS resource
COPY --from=builder /assets /opt/resource

FROM resource
