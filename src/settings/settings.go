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
	WebserverAddr = ":8100"

	// Postgres
	DatabaseConnRetries = 5
	DatabaseAddr        = "honeypot-postgres:5432"
	DatabaseDatabase    = "honeypot_db"
	DatabaseUser        = "honeypot_user"
	DatabasePassword    = "honeypot_pass"

	// Rabbitmq
	RabbitmqHost                  = "honeypot-rabbitmq"
	RabbitmqPort                  = "5672"
	RabbitmqTaskProcessConnAttemp = "ProcessConnAttemp"
	RabbitmqTaskStoreConnAttemp   = "StoreConnAttemp"

	// Connection Retry
	ConnectionRetriesWait  = time.Second * 5
	ConnectionRetriesTotal = 10

	// Page defaults
	PageDefaultNum  = 0
	PageDefaultSize = 10
)
