package timelines

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ConnAttemp struct {
	Time        time.Time
	Port        string
	IP          string
	CountryCode string
	ClientPort  string
}

type MapDataEntry struct {
	CountryCode string
	Count       int64
}

func NewConnAttemp(tm time.Time, port, addr string) (*ConnAttemp, error) {
	ipAndPort, err := separateIPAndPort(addr)
	if err != nil {
		return nil, err
	}

	connAttem := &ConnAttemp{
		Time:        tm,
		Port:        port,
		IP:          ipAndPort[0],
		CountryCode: "",
		ClientPort:  ipAndPort[1],
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

	return rep
}

func (c *ConnAttemp) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func ConnAttempFromJson(b []byte) (*ConnAttemp, error) {
	var connAttemp *ConnAttemp
	err := json.Unmarshal(b, &connAttemp)
	if err != nil {
		return nil, err
	}
	return connAttemp, nil
}

func separateIPAndPort(addr string) ([]string, error) {
	strs := strings.Split(addr, ":")
	if len(strs) != 2 {
		return nil, fmt.Errorf("Addr can't be separated by ':', %s", addr)
	}
	return strs, nil
}
