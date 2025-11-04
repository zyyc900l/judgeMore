FROM golang:1.25 AS builder

ENV TZ=Asia/Shanghai
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# 复制所有文件到容器中
COPY . .

RUN go mod tidy
RUN sh build.sh

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
ENV service=video

WORKDIR /app
COPY --from=builder /app/output /app/output
COPY --from=builder /app/config /app/config

CMD ["sh", "./output/bootstrap.sh"]