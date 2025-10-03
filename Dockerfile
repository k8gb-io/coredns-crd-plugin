FROM alpine:3.22.1
COPY coredns /
EXPOSE 53 53/udp
ENTRYPOINT ["/coredns"]
