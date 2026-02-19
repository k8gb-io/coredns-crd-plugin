FROM gcr.io/distroless/static-debian12:nonroot
COPY coredns /
EXPOSE 53 53/udp
ENTRYPOINT ["/coredns"]
