package timelines

import (
	"context"
	"fmt"
	"honeypot/settings"
	"log"

	"github.com/influxdata/influxdb-client-go/v2/api"
)

type TimelinesQuery interface {
	Close()
	GetConnAttemps(context.Context) ([]*ConnAttemp, error)
}

type timelinesQuery struct {
	*timelines
	queryAPI api.QueryAPI
}

func NewTimelinesQuery() TimelinesQuery {
	log.Println("Timelines starting")

	tl := newTimelines()
	queryAPI := tl.dbClient.QueryAPI(settings.InfluxDBOrg)

	return &timelinesQuery{
		timelines: tl,
		queryAPI:  queryAPI,
	}
}

func (t *timelinesQuery) Close() {
	t.timelines.close()
	log.Println("Timelines closed")
}

func (t *timelinesQuery) GetConnAttemps(ctx context.Context) ([]*ConnAttemp, error) {
	query := makeConnAttempsQuery()
	result, err := t.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	ret := make([]*ConnAttemp, 0)
	for result.Next() {
		record := result.Record()

		port, ok := record.ValueByKey("Port").(string)
		if !ok {
			return nil, fmt.Errorf("Port is not a string")
		}
		ip, ok := record.ValueByKey("IP").(string)
		if !ok {
			return nil, fmt.Errorf("IP is not a string")
		}
		clientPort, ok := record.ValueByKey("ClientPort").(string)
		if !ok {
			clientPort = ""
		}
		countryCode, ok := record.ValueByKey("CountryCode").(string)
		if !ok {
			countryCode = ""
		}

		connAttemp := &ConnAttemp{
			Time:        record.Time(),
			Port:        port,
			IP:          ip,
			ClientPort:  clientPort,
			CountryCode: countryCode,
		}
		ret = append(ret, connAttemp)
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	return ret, nil
}

func makeConnAttempsQuery() string {
	return `from(bucket:"honeypot")
		|> range(start: -10h)
		|> filter(fn: (r) => r._measurement == "conn")
		|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")`
}
