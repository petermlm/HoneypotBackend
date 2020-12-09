package timelines

import (
	"honeypot/settings"
	"log"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type TimelinesWiter interface {
	Close()
	InsertConnAttemp(connAttemp *ConnAttemp)
}

type timelinesWriter struct {
	*timelines
	writeAPI api.WriteAPI
}

func NewTimelinesWriter() TimelinesWiter {
	log.Println("Timelines starting")

	tl := newTimelines()
	writeAPI := tl.dbClient.WriteAPI(settings.InfluxDBOrg, settings.InfluxDBBucket)

	return &timelinesWriter{
		timelines: tl,
		writeAPI:  writeAPI,
	}
}

func (t *timelinesWriter) Close() {
	t.writeAPI.Flush()
	t.timelines.close()
	log.Println("Timelines closed")
}

func (t *timelinesWriter) InsertConnAttemp(connAttemp *ConnAttemp) {
	point := connAttemp.toDbPoint()
	t.insertCommonPoint("conn", point)
}

func (t *timelinesWriter) insertCommonPoint(measure string, dbp *dbPoint) {
	point := influxdb2.NewPoint(measure, dbp.Tags, dbp.Fields, dbp.Timestamp)
	t.writeAPI.WritePoint(point)
}
