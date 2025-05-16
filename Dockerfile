FROM alpine:3.21.3
COPY coredns /
EXPOSE 53 53/udp
ENTRYPOINT ["/coredns"]
