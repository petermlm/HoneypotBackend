package timelines

import (
	"context"
	"fmt"
	"honeypot/settings"
	"log"
	"strings"

	"github.com/influxdata/influxdb-client-go/v2/api"
)

type TimelinesQuery interface {
	Close()
	GetTotalConsumptions(context.Context, string) (*SingleCount, error)
	GetMapData(context.Context, string) ([]*MapDataEntry, error)
	GetConnAttemps(context.Context, string) ([]*ConnAttemp, error)
	GetTopConsumers(context.Context, string) ([]*MapDataEntry, error)
	GetTopFlavours(context.Context, string) ([]*PortCount, error)
	GetBytes(context.Context, string) ([]*BytesList, error)
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

func (t *timelinesQuery) getCommon(
	ctx context.Context,
	rangeValue string,
	fluxQuery func(string) (string, error),
) (*api.QueryTableResult, error) {
	query, err := fluxQuery(rangeValue)
	if err != nil {
		return nil, err
	}

	result, err := t.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *timelinesQuery) GetTotalConsumptions(ctx context.Context, rangeValue string) (*SingleCount, error) {
	result, err := t.getCommon(ctx, rangeValue, makeTotalConsumptions)
	if err != nil {
		return nil, err
	}

	if !result.Next() {
		return nil, fmt.Errorf("No results")
	}

	count, ok := result.Record().ValueByKey("Count").(int64)
	if !ok {
		return nil, fmt.Errorf("No results")
	}
	return &SingleCount{int(count)}, nil
}

func (t *timelinesQuery) GetMapData(ctx context.Context, rangeValue string) ([]*MapDataEntry, error) {
	result, err := t.getCommon(ctx, rangeValue, makeMapDataQuery)
	if err != nil {
		return nil, err
	}

	ret := make([]*MapDataEntry, 0)
	for result.Next() {
		record := result.Record()

		countryCode, ok := record.ValueByKey("CountryCode").(string)
		if !ok {
			countryCode = ""
		}
		count, ok := record.Value().(int64)
		if !ok {
			log.Println(record.Value())
			count = 0
		}

		mapDataEntry := &MapDataEntry{
			CountryCode: countryCode,
			Count:       count,
		}
		ret = append(ret, mapDataEntry)
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	return ret, nil
}

func (t *timelinesQuery) GetConnAttemps(ctx context.Context, rangeValue string) ([]*ConnAttemp, error) {
	result, err := t.getCommon(ctx, rangeValue, makeConnAttempsQuery)
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

func (t *timelinesQuery) GetTopConsumers(ctx context.Context, rangeValue string) ([]*MapDataEntry, error) {
	result, err := t.getCommon(ctx, rangeValue, makeTopConsumersQuery)
	if err != nil {
		return nil, err
	}

	ret := make([]*MapDataEntry, 0)
	for result.Next() {
		record := result.Record()

		countryCode, ok := record.ValueByKey("CountryCode").(string)
		if !ok {
			countryCode = ""
		}
		count, ok := record.Value().(int64)
		if !ok {
			log.Println(record.Value())
			count = 0
		}

		portCount := &MapDataEntry{
			CountryCode: countryCode,
			Count:       count,
		}
		ret = append(ret, portCount)
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	return ret, nil
}

func (t *timelinesQuery) GetTopFlavours(ctx context.Context, rangeValue string) ([]*PortCount, error) {
	result, err := t.getCommon(ctx, rangeValue, makeTopFlavoursQuery)
	if err != nil {
		return nil, err
	}

	ret := make([]*PortCount, 0)
	for result.Next() {
		record := result.Record()

		port, ok := record.ValueByKey("Port").(string)
		if !ok {
			port = ""
		}
		count, ok := record.Value().(int64)
		if !ok {
			log.Println(record.Value())
			count = 0
		}

		portCount := &PortCount{
			Port:  port,
			Count: count,
		}
		ret = append(ret, portCount)
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	return ret, nil
}

func (t *timelinesQuery) GetBytes(ctx context.Context, rangeValue string) ([]*BytesList, error) {
	result, err := t.getCommon(ctx, rangeValue, makeBytesElasticsearchQuery)
	if err != nil {
		return nil, err
	}

	ret := make([]*BytesList, 0)
	for result.Next() {
		record := result.Record()

		bytes, ok := record.ValueByKey("Bytes").(string)
		if !ok {
			bytes = ""
		}
		bytes = strings.ReplaceAll(bytes, "\u0000", "")
		bytes = strings.TrimSpace(bytes)

		byteList := &BytesList{
			Time:  record.Time(),
			Bytes: bytes,
		}
		ret = append(ret, byteList)
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	return ret, nil
}
