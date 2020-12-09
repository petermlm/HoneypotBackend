package timelines

import (
	"fmt"
	"honeypot/settings"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type timelines struct {
	dbClient influxdb2.Client
}

func (t *timelines) close() {
	t.dbClient.Close()
}

func newTimelines() *timelines {
	addr := makeInfluxDBAddr()
	token := makeInfluxDBToken()
	dbClient := influxdb2.NewClient(addr, token)
	return &timelines{dbClient}
}

func makeInfluxDBAddr() string {
	return fmt.Sprintf("http://%s:%s", settings.InfluxDBAddr, settings.InfluxDBPort)
}

func makeInfluxDBToken() string {
	return fmt.Sprintf("%s:%s", settings.InfluxDBUsername, settings.InfluxDBPassword)
}
