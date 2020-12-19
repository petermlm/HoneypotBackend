package timelines

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRangeFluxFunc_GoodCall(t *testing.T) {
	for key, val := range validRanges {
		rangeAct, err := makeRangeFluxFunc(key)
		rangeExp := fmt.Sprintf("|> range(start: -1%s)", key)
		assert.True(t, val)
		assert.NoError(t, err)
		assert.Equal(t, rangeAct, rangeExp)
	}
}

func TestMakeRangeFluxFunc_InvalidRange(t *testing.T) {
	invalivRange := "invalid"
	_, errAct := makeRangeFluxFunc(invalivRange)
	errExp := newInvalidRange(invalivRange)
	assert.EqualError(t, errAct, errExp.Error())
}

func TestQueryCommon_ValidRange(t *testing.T) {
	checkPart := func(query, part string) {
		assert.NotEqual(t, strings.Index(query, part), -1)
	}

	rangeValue := "mo"
	rangePart, err := makeRangeFluxFunc(rangeValue)
	assert.NoError(t, err)

	queryAct, err := makeQueryCommon(rangeValue, "QUERY")
	assert.NoError(t, err)

	checkPart(queryAct, `from(bucket:"honeypot")`)
	checkPart(queryAct, rangePart)
	checkPart(queryAct, `QUERY`)
}

func TestQueryCommon_InvalidRange(t *testing.T) {
	invalivRange := "invalid"
	_, errAct := makeQueryCommon(invalivRange, "QUERY")
	errExp := newInvalidRange(invalivRange)
	assert.EqualError(t, errAct, errExp.Error())
}
