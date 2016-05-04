FROM ubuntu
COPY video-provisioner /bin/video-provisioner
ENTRYPOINT ["video-provisioner"]
CMD ["--help"]