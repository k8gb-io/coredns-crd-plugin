FROM alpine
COPY coredns /
EXPOSE 53 53/udp
ENTRYPOINT ["/coredns"]
