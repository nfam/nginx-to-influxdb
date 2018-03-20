package main

import (
    "fmt"
    "bytes"
    "net/http"
    "os"
    "strings"
    "gopkg.in/mcuadros/go-syslog.v2"
)

func main() {

    // Example
    // url := "http://10.0.0.2:8086/write?db=nginx&precision=ms"

    url := os.Getenv("INFLUXDB_WRITE_URL")
    if url == "" {
        fmt.Println("Environment INFLUXDB_WRITE_URL is not provided.")
        return
    }

    channel := make(syslog.LogPartsChannel)
    handler := syslog.NewChannelHandler(channel)

    server := syslog.NewServer()
    server.SetFormat(syslog.RFC3164)
    server.SetHandler(handler)
    server.ListenUDP("0.0.0.0:514")
    server.ListenTCP("0.0.0.0:514")

    server.Boot()

    go func(channel syslog.LogPartsChannel) {
        for logParts := range channel {
            content := fixTimestamp(logParts["content"].(string))
            body := []byte(content)
            req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "text/plain")
            client := &http.Client{}
            resp, err := client.Do(req)
            if err != nil {
                panic(err)
            }
            resp.Body.Close()
        }
    } (channel)

    server.Wait()
}

// Converts nginx $msec (Float) to milliseconds (Int).
func fixTimestamp(content string) string {
    timestampIndex := strings.LastIndex(content, " ")
    timestampParts := strings.Split(content[timestampIndex + 1:], ".")
    timestamp := timestampParts[0]
    if len(timestampParts) > 1 {
        timestamp += (timestampParts[1] + "000")[0:3]
    } else {
        timestamp += "000"
    }

    return content[:timestampIndex] + " " + timestamp
}
