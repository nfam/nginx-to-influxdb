# build
FROM golang:1.9.4-alpine3.7
RUN mkdir /nginx-to-influxdb
WORKDIR /nginx-to-influxdb
COPY main.go /nginx-to-influxdb/
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN apk add --no-cache git && \
    go get gopkg.in/mcuadros/go-syslog.v2 && \
    go build  -ldflags '-w -s' -a -installsuffix cgo -o nginx-to-influxdb

# image
FROM scratch
COPY --from=0 /nginx-to-influxdb/nginx-to-influxdb /
EXPOSE 514 514//udp
ENTRYPOINT ["/nginx-to-influxdb"]