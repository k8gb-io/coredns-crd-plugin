FROM golang:1.25.7-alpine
COPY coredns /
EXPOSE 53 53/udp
ENTRYPOINT ["/coredns"]
