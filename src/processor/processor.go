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
		resp, err := http.Get(makeGeoURL(m.IP))
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		var objmap map[string]json.RawMessage
		if err := json.Unmarshal(body, &objmap); err != nil {
			continue
		}
		m.CountryCode = string(objmap["geoplugin_countryCode"])[1:3]

		p.Publish(m)
	}

	return nil
}
