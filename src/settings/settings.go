package settings

// Ports to be listent to
var Ports = [...]string{
	"5432",
	"3306",
	"27017",
	"7474",
	"9200",
}

const (
	// Webserver
	WebserverAddr = ":8080"

	// InfluxDB
	InfluxDBAddr      = "localhost"
	InfluxDBPort      = "8086"
	InfluxDBUsername  = "honey"
	InfluxDBPassword  = "honey"
	InfluxDBName      = "honey"
	InfluxDBPrecision = "ns"
	InfluxDBOrg       = "honeypot"
	InfluxDBBucket    = "honeypot"

	// Rabbitmq
	RabbitmqHost                  = "localhost"
	RabbitmqPort                  = "5672"
	RabbitmqTaskProcessConnAttemp = "ProcessConnAttemp"
	RabbitmqTaskStoreConnAttemp   = "StoreConnAttemp"

	// Page defaults
	PageDefaultNum  = 0
	PageDefaultSize = 10
)
