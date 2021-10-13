FROM alpine:3.14.2
COPY coredns /
EXPOSE 53 53/udp
ENTRYPOINT ["/coredns"]
