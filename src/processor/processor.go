package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"honeypot/queue"
	"honeypot/settings"
	"honeypot/timelines"
	"io/ioutil"
	"net/http"
)

type env struct {
	c  queue.Consumer
	tl timelines.TimelinesWriter
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

	// e := env{
	// 	c:  c,
	// 	tl: tl,
	// }

	ch, _ := c.Consume(context.Background())
	for m := range ch {
		resp, err := http.Get(fmt.Sprintf("http://www.geoplugin.net/json.gp?ip=%s", m.IP))
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
