FROM golang:1.5.2

COPY . /go/src/github.com/elleFlorio/testAppGru
WORKDIR /go/src/github.com/elleFlorio/testAppGru

ENV GOPATH /go/src/github.com/elleFlorio/testAppGru:$GOPATH
RUN CGO_ENABLED=0 go install github.com/elleFlorio/testAppGru

ENTRYPOINT ["testAppGru"]
CMD ["--help"]