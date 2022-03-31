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

type DBObject struct {
	ID        uint
	CreatedAt time.Time `sql:"default:now()"`
	UpdatedAt time.Time `sql:"default:now()"`
}

type ConnAttemp struct {
	DBObject

	Time        time.Time `sql:",notnull"`
	Port        string    `sql:",notnull"`
	IP          string    `sql:",notnull"`
	CountryCode string    `sql:",notnull"`
	ClientPort  string    `sql:",notnull"`
	Bytes       []byte    `sql:",notnull"`
}

func NewConnAttemp(tm time.Time, port, addr string) (*ConnAttemp, error) {
	ipAndPort, err := separateIPAndPort(addr)
	if err != nil {
		return nil, err
	}
	ip := ipAndPort[0]
	clientPort := ipAndPort[1]

	if !reIPv4.Match([]byte(ip)) {
		return nil, newInvalidIP(ip)
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

func separateIPAndPort(addr string) ([]string, error) {
	strs := strings.Split(addr, ":")
	if len(strs) != 2 {
		return nil, newCantSeparateAddrError(addr)
	}
	return strs, nil
}

/* ============================================================================
 * Move to query models
 * ============================================================================
 */

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
