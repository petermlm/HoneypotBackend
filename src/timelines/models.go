package timelines

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"
)

const ipv4Str = `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`

var reIPv4 *regexp.Regexp

func init() {
	reIPv4 = regexp.MustCompile(ipv4Str)
}

type ConnAttemp struct {
	Time        time.Time
	Port        string
	IP          string
	CountryCode string
	ClientPort  string
	Bytes       []byte
}

type MapDataEntry struct {
	CountryCode string
	Count       int64
}

type PortCount struct {
	Port  string
	Count int64
}

func NewConnAttemp(tm time.Time, port, addr string) (*ConnAttemp, error) {
	ipAndPort, err := separateIPAndPort(addr)
	if err != nil {
		return nil, err
	}
	ip := ipAndPort[0]
	clientPort := ipAndPort[1]

	if !reIPv4.Match([]byte(ip)) {
		return nil, newInvalidIP7(ip)
	}

	connAttem := &ConnAttemp{
		Time:        tm,
		Port:        port,
		IP:          ip,
		CountryCode: "",
		ClientPort:  clientPort,
		Bytes:       make([]byte, 0),
	}
	return connAttem, nil
}

type dbPoint struct {
	Timestamp time.Time
	Tags      map[string]string
	Fields    map[string]interface{}
}

func (c *ConnAttemp) toDbPoint() *dbPoint {
	rep := new(dbPoint)

	rep.Timestamp = c.Time

	rep.Tags = make(map[string]string)
	rep.Tags["Port"] = c.Port
	rep.Tags["IP"] = c.IP
	rep.Tags["CountryCode"] = c.CountryCode

	rep.Fields = make(map[string]interface{})
	rep.Fields["ClientPort"] = c.ClientPort
	rep.Fields["Bytes"] = c.Bytes

	return rep
}

func (c *ConnAttemp) ToJSON() ([]byte, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return nil, newCantMarshalConnAttemp(err)
	}
	return b, nil
}

func ConnAttempFromJson(b []byte) (*ConnAttemp, error) {
	var connAttemp *ConnAttemp
	err := json.Unmarshal(b, &connAttemp)
	if err != nil {
		return nil, newCantUnarshalConnAttemp(err)
	}
	return connAttemp, nil
}

func separateIPAndPort(addr string) ([]string, error) {
	strs := strings.Split(addr, ":")
	if len(strs) != 2 {
		return nil, newCantSeparateAddrError(addr)
	}
	return strs, nil
}
