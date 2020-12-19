package timelines

import (
	"encoding/base64"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NewConnAttempSuite struct {
	suite.Suite
	timeObj    time.Time
	port       string
	ip         string
	clientPort string
	addr       string
}

func (s *NewConnAttempSuite) SetupTest() {
	var err error

	timeStr := "2020-12-18T20:02:23.049563482+01:00"
	s.timeObj, err = time.Parse(time.RFC3339, timeStr)
	s.port = "9000"
	s.ip = "123.123.123.123"
	s.clientPort = "98765"
	s.addr = fmt.Sprintf("%s:%s", s.ip, s.clientPort)

	assert.NoError(s.T(), err)
}

func (s *NewConnAttempSuite) TestNewConnAttemp_GoodCall() {
	connAttemp, err := NewConnAttemp(s.timeObj, s.port, s.addr)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), connAttemp.Time.String(), s.timeObj.String())
	assert.Equal(s.T(), connAttemp.Port, s.port)
	assert.Equal(s.T(), connAttemp.IP, s.ip)
	assert.Equal(s.T(), connAttemp.CountryCode, "")
	assert.Equal(s.T(), connAttemp.ClientPort, s.clientPort)
	assert.Equal(s.T(), len(connAttemp.Bytes), 0)
}

func (s *NewConnAttempSuite) TestNewConnAttemp_IPMultipleColumns() {
	addr := "123.123.123:123:123"
	_, errAct := NewConnAttemp(s.timeObj, s.port, addr)
	errExp := newCantSeparateAddrError(addr)
	assert.EqualError(s.T(), errAct, errExp.Error())
}

func (s *NewConnAttempSuite) TestNewConnAttemp_IPInvalidCall() {
	ip := "123.123.123"
	port := "9000"
	addr := fmt.Sprintf("%s:%s", ip, port)
	_, errAct := NewConnAttemp(s.timeObj, s.port, addr)
	errExp := newInvalidIP7(ip)
	assert.EqualError(s.T(), errAct, errExp.Error())
}

func TestNewConnAttempSuite(t *testing.T) {
	suite.Run(t, new(NewConnAttempSuite))
}

type ConnAttempJSONSuite struct {
	suite.Suite
	connAttemp *ConnAttemp
}

func (s *ConnAttempJSONSuite) SetupTest() {
	s.connAttemp, _ = NewConnAttemp(time.Now(), "123", "123.123.123.123:987")
}

func (s *ConnAttempJSONSuite) TestConnAttempJSON_MarshalValid() {
	// s.connAttemp.Bytes = bytesB64
	_, err := s.connAttemp.ToJSON()
	assert.NoError(s.T(), err)
}

// TODO: How to test this?
// func (s *ConnAttempJSONSuite) TestConnAttempJSON_MarshalInvalid() {
// }

func (s *ConnAttempJSONSuite) TestConnAttempJSON_UnmarshalValid() {
	timeStr := "2020-12-18T20:02:23.049563482+01:00"
	port := "123"
	ip := "123.123.123.123"
	countryCode := "PT"
	clientPort := "987"
	bytes := []byte("experiment")
	bytesB64 := base64.StdEncoding.EncodeToString(bytes)

	json := []byte(fmt.Sprintf(`{
		"Time":"%s",
		"Port":"%s",
		"IP":"%s",
		"CountryCode":"%s",
		"ClientPort":"%s",
		"Bytes":"%s"
	}`, timeStr, port, ip, countryCode, clientPort, bytesB64))
	timeExp, err := time.Parse(time.RFC3339, timeStr)
	assert.NoError(s.T(), err)

	connAttemp, err := ConnAttempFromJson(json)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), connAttemp.Time.String(), timeExp.String())
	assert.Equal(s.T(), connAttemp.Port, port)
	assert.Equal(s.T(), connAttemp.IP, ip)
	assert.Equal(s.T(), connAttemp.CountryCode, countryCode)
	assert.Equal(s.T(), connAttemp.ClientPort, clientPort)
	assert.Equal(s.T(), connAttemp.Bytes, bytes)
}

func (s *ConnAttempJSONSuite) TestConnAttempJSON_UnmarshalInvalid() {
	json := "not a json"
	_, err := ConnAttempFromJson([]byte(json))
	var errExp *CantUnmarshalConnAttemp
	if !errors.As(err, &errExp) {
		s.T().Fatalf("unexpected error: %v", err)
	}
}

func TestConnAttempJSONSuite(t *testing.T) {
	suite.Run(t, new(ConnAttempJSONSuite))
}

func TestSeparateIPAndPort_GoodCall(t *testing.T) {
	ip := "123.123.123.123"
	port := "9000"
	addr := fmt.Sprintf("%s:%s", ip, port)
	act, err := separateIPAndPort(addr)
	assert.Equal(t, len(act), 2)
	assert.Equal(t, act[0], ip)
	assert.Equal(t, act[1], port)
	assert.NoError(t, err)
}

func TestSeparateIPAndPort_MultipleColumns(t *testing.T) {
	addr := "123.123.123:123:123"
	_, errAct := separateIPAndPort(addr)
	errExp := newCantSeparateAddrError(addr)
	assert.EqualError(t, errAct, errExp.Error())
}

func TestSeparateIPAndPort_InvalidCall(t *testing.T) {
	addr := "not valid"
	_, errAct := separateIPAndPort(addr)
	errExp := newCantSeparateAddrError(addr)
	assert.EqualError(t, errAct, errExp.Error())
}
