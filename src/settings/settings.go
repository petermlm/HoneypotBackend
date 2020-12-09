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
	// InfluxDB connection settings
	InfluxDBAddr      = "localhost"
	InfluxDBPort      = "8086"
	InfluxDBUsername  = "honey"
	InfluxDBPassword  = "honey"
	InfluxDBName      = "honey"
	InfluxDBPrecision = "ns"
	InfluxDBOrg       = "honeypot"
	InfluxDBBucket    = "honeypot"
)
