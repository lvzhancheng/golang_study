FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn \
    GOARCH=amd64 \

    WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN env
RUN go mod download

COPY 2.2/ .

RUN go build -o http_server .

FROM alpine

COPY --from=builder /build/http_server /

ENTRYPOINT ["/http_server"]