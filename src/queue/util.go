package queue

import (
	"fmt"
	"honeypot/settings"
)

func makeConnString() string {
	return fmt.Sprintf("amqp://guest:guest@%s:%s/", settings.RabbitmqHost, settings.RabbitmqPort)
}
