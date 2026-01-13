FROM alpine:3.23.2
COPY coredns /
EXPOSE 53 53/udp
ENTRYPOINT ["/coredns"]
