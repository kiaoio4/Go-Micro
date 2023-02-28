FROM golang:1.18-alpine as builder
COPY . .
RUN bash build.sh go-micro /go-micro

COPY go-micro.toml /go-micro.toml
COPY swagger /swagger
WORKDIR /
EXPOSE  80
HEALTHCHECK --interval=30s --timeout=15s \
    CMD curl --fail http://localhost:80/health || exit 1
COPY --from=builder /go-micro /go-micro
ENTRYPOINT [ "/go-micro" ]
CMD ["run"]
