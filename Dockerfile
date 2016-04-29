FROM golang:1.5.4-alpine

COPY . /go/src/github.com/elleFlorio/video-provisioner
WORKDIR /go/src/github.com/elleFlorio/video-provisioner

ENV GOPATH /go/src/github.com/elleFlorio/video-provisioner:$GOPATH
RUN CGO_ENABLED=0 go install github.com/elleFlorio/video-provisioner

RUN CGO_ENABLED=0 go install github.com/elleFlorio/video-provisioner

ENTRYPOINT ["video-provisioner"]
CMD ["--help"]