FROM golang:1.11 as builder
ARG VERSION
ARG GIT_COMMIT
WORKDIR /go/src/github.com/platanus/nchan-prometheus-exporter
COPY *.go ./
COPY vendor ./vendor
COPY collector ./collector
COPY nginxClient ./nginxClient
COPY nchanClient ./nchanClient
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT}" -o exporter .

FROM alpine:latest
COPY --from=builder /go/src/github.com/platanus/nchan-prometheus-exporter/exporter /usr/bin/
ENTRYPOINT [ "/usr/bin/exporter" ]