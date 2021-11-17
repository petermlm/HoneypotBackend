package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"honeypot/queue"
	"honeypot/settings"
	"io/ioutil"
	"net/http"
)

const baseGeoURL = "http://www.geoplugin.net/json.gp"

func makeGeoURL(ip string) string {
	return fmt.Sprintf("%s?ip=%s", baseGeoURL, ip)
}

func Start() error {
	c, err := queue.NewConsumer(settings.RabbitmqTaskProcessConnAttemp)
	if err != nil {
		return err
	}
	defer c.Destroy()

	p, err := queue.NewPublisher(settings.RabbitmqTaskStoreConnAttemp)
	if err != nil {
		return err
	}
	defer p.Destroy()

	ch, _ := c.Consume(context.Background())
	for m := range ch {
		m.CountryCode = getCountryCode(m.IP)
		p.Publish(m)
	}

	return nil
}

func getCountryCode(IP string) string {
	resp, err := http.Get(makeGeoURL(IP))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var objmap map[string]json.RawMessage
	if err := json.Unmarshal(body, &objmap); err != nil {
		return ""
	}
	countryCode := string(objmap["geoplugin_countryCode"])
	return countryCode[1:3]
}
