FROM golang:1.20.2-alpine3.17 as builder

RUN apk add --no-cache git make

ENV GOPROXY=https://goproxy.cn
ENV SRC_DIR=/go/src/goutils

RUN git clone --depth 1 https://github.com/lonelypale/goutils.git $SRC_DIR

RUN cd $SRC_DIR \
    && go build -ldflags="-s -w" -installsuffix cgo -o target/sysinfod cmd/sysinfo/main.go

FROM alpine:3.17

ENV TZ Asia/Shanghai
RUN apk add --no-cache ca-certificates tzdata \
    && cp /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && apk del tzdata

COPY --from=builder /go/src/goutils/target/sysinfod /usr/local/bin/sysinfod

ENV DATA_PATH=/data
RUN mkdir -p $DATA_PATH
VOLUME $DATA_PATH
WORKDIR $DATA_PATH

EXPOSE 9999

CMD ["sysinfod"]
