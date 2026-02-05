FROM alpine:3.23.3
COPY coredns /
EXPOSE 53 53/udp
ENTRYPOINT ["/coredns"]
