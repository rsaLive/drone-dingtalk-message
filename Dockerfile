FROM golang:alpine as builder

LABEL stage=gobuilder
ENV TZ Asia/Shanghai
ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /build
COPY ./build/drone-ding /app/drone-ding

# run step
FROM alpine

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tzdata

WORKDIR /app
# copy bin from build step
COPY --from=builder /app/drone-ding .

ENTRYPOINT ["/app/drone-ding"]