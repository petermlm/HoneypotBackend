package settings

import "time"

// Ports to be listent to
var Ports = [...]string{
	"3306",  // MySQL
	"5432",  // PostgreSQL
	"7474",  // Neo4j
	"9200",  // Elasticsearch
	"27017", // MongoDB
}

const (
	// Webserver
	WebserverAddr = ":8080"

	// InfluxDB
	InfluxDBAddr      = "influxdb"
	InfluxDBPort      = "8086"
	InfluxDBUsername  = "honey"
	InfluxDBPassword  = "honey"
	InfluxDBName      = "honey"
	InfluxDBPrecision = "ns"
	InfluxDBOrg       = "honeypot"
	InfluxDBBucket    = "honeypot"

	// Rabbitmq
	RabbitmqHost                  = "rabbitmq"
	RabbitmqPort                  = "5672"
	RabbitmqTaskProcessConnAttemp = "ProcessConnAttemp"
	RabbitmqTaskStoreConnAttemp   = "StoreConnAttemp"

	// Connection Retry
	ConnectionRetriesWait  = time.Second * 2
	ConnectionRetriesTotal = 10

	// Page defaults
	PageDefaultNum  = 0
	PageDefaultSize = 10
)
