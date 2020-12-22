package timelines

import "fmt"

const DefaultValidRanges string = "mo"

var validRanges = map[string]bool{
	"y":  true,
	"mo": true,
	"w":  true,
	"d":  true,
	"h":  true,
}

func makeRangeFluxFunc(rangeValue string) (string, error) {
	if _, ok := validRanges[rangeValue]; !ok {
		return "", newInvalidRange(rangeValue)
	}
	rangeStr := fmt.Sprintf("|> range(start: -1%s)", rangeValue)
	return rangeStr, nil
}

func makeQueryCommon(rangeValue, queryPart string) (string, error) {
	rangeFluxFunc, err := makeRangeFluxFunc(rangeValue)
	if err != nil {
		return "", err
	}
	queryCommon := fmt.Sprintf(`from(bucket:"honeypot") %s %s`, rangeFluxFunc, queryPart)
	return queryCommon, nil
}

func makeTotalConsumptions(rangeValue string) (string, error) {
	return makeQueryCommon(rangeValue, `
		|> pivot(rowKey: ["_time", "IP", "CountryCode", "Port"], columnKey: ["_field"], valueColumn: "_value")
        |> group()
  		|> count(column: "IP")
		|> rename(columns: {"IP": "Count"})`)
}

func makeMapDataQuery(rangeValue string) (string, error) {
	return makeQueryCommon(rangeValue, `
		|> filter(fn: (r) => r._measurement == "conn")
		|> pivot(rowKey: ["_time", "IP", "CountryCode", "Port"], columnKey: ["_field"], valueColumn: "_value")
		|> group(columns: ["CountryCode"], mode:"by")
		|> count(column: "IP")
		|> rename(columns: {"IP": "_value"})`)
}

func makeConnAttempsQuery(rangeValue string) (string, error) {
	return makeQueryCommon(rangeValue, `
		|> filter(fn: (r) => r._measurement == "conn")
		|> pivot(rowKey: ["_time", "IP", "CountryCode", "Port"], columnKey: ["_field"], valueColumn: "_value")
        |> drop(columns: ["Bytes", "ClientPort"])
		|> group()
		|> sort(columns: ["_time"], desc: true)`)
}

func makeTopConsumersQuery(rangeValue string) (string, error) {
	return makeQueryCommon(rangeValue, `
		|> filter(fn: (r) => r._measurement == "conn")
		|> pivot(rowKey: ["_time", "IP", "CountryCode", "Port"], columnKey: ["_field"], valueColumn: "_value")
		|> group(columns: ["CountryCode"], mode:"by")
		|> count(column: "IP")
        |> group()
        |> sort(columns: ["IP"], desc: true)
  		|> limit(n: 10, offset: 0)
		|> rename(columns: {"IP": "_value"})`)
}

func makeTopFlavoursQuery(rangeValue string) (string, error) {
	return makeQueryCommon(rangeValue, `
		|> filter(fn: (r) => r._measurement == "conn")
		|> pivot(rowKey: ["_time", "IP", "CountryCode", "Port"], columnKey: ["_field"], valueColumn: "_value")
		|> group(columns: ["Port"], mode:"by")
        |> count(column: "IP")
        |> group()
        |> sort(columns: ["IP"], desc: true)
  		|> limit(n: 10, offset: 0)
		|> rename(columns: {"IP": "_value"})`)
}
