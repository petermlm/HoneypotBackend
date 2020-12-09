package timelines

import (
	"fmt"
	"strings"
	"time"
)

type ConnAttemp struct {
	Time time.Time
	Port string
	Addr string
}

type dbPoint struct {
	Timestamp time.Time
	Tags      map[string]string
	Fields    map[string]interface{}
}

func (c *ConnAttemp) toDbPoint() *dbPoint {
	rep := new(dbPoint)

	ipAndPort, err := separateIPAndPort(c.Addr)
	if err != nil {
		// TODO: Handle
	}

	rep.Timestamp = c.Time

	rep.Tags = make(map[string]string)
	rep.Tags["Port"] = c.Port
	rep.Tags["Addr"] = ipAndPort[0]

	rep.Fields = make(map[string]interface{})
	rep.Fields["ClientPort"] = ipAndPort[1]

	return rep
}

func separateIPAndPort(addr string) ([]string, error) {
	strs := strings.Split(addr, ":")
	if len(strs) != 2 {
		return nil, fmt.Errorf("Addr can't be separated by ':', %s", addr)
	}
	return strs, nil
}
