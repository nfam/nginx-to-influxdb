# nginx-to-influxdb

A small syslog server to forward nginx log to InfluxDB.

## Set up `nginx-to-influxdb` syslog server

```sh
docker run -d \
    --name=nginx-to-influxdb \
    -p 514:514 \
    -p 514:514/udp \
    -e INFLUXDB_WRITE_URL="http://${HOST}/write?db=${DATABASE}&precision=ms" \
    nfam/nginx-to-influxdb
```

*NOTE*: `INFLUXDB_WRITE_URL` must have `precision=ms` in its query string as `Nginx` supports time resolution up to milliseconds only.

## Configure `Nginx`

Nginx log format must comply InfluxDB [Line Protocol](https://docs.influxdata.com/influxdb/v1.5/write_protocols/line_protocol_tutorial/)

```conf
log_format      access_line_protocol 'nginx'
                ' body_bytes_sent=$body_bytes_sent'
                ',bytes_sent=$bytes_sent'
                ',host="$host"'
                ',http_referrer="$http_referer"'
                ',http_user_agent="$http_user_agent"'
                ',remote_addr="$remote_addr"'
                ',request_length=$request_length'
                ',request_method="$request_method"'
                ',request_time=$request_time'
                ',request_uri="$request_uri"'
                ',server_protocol="$server_protocol"'
                ',status=$status'
                ',upstream_response_time="$upstream_response_time"'
                ' '
                '$msec';
access_log      syslog:server=${nginx-to-influxdb} access_line_protocol;
```
