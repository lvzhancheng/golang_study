FROM golang:1.18 AS builder
# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn \
    GOARCH=amd64
WORKDIR /build
COPY ./ .
RUN env
RUN go mod download
RUN go build -o http_server .

FROM alpine:3.10
COPY --from=builder /build/http_server /
COPY config.toml /
ENTRYPOINT ["/http_server"]